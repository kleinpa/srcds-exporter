"""Prometheus exporter for Source Dedicated Server (srcds) via RCON."""

import argparse
import logging
import os
import re

from prometheus_client import start_http_server
from prometheus_client.core import CollectorRegistry, GaugeMetricFamily, InfoMetricFamily
from prometheus_client.registry import Collector
from rcon.source import Client

log = logging.getLogger("srcds_exporter")

# `stats` header columns mapped to (metric suffix, help). CS:GO and older Source
# games (e.g. CS:Source) label and order the columns differently.
STATS_COLUMNS = {
    "CPU": ("cpu", "CPU usage percent"),
    "NetIn": ("netin", "Incoming network traffic"),
    "In_(KB/s)": ("netin", "Incoming network traffic"),
    "NetOut": ("netout", "Outgoing network traffic"),
    "Out_(KB/s)": ("netout", "Outgoing network traffic"),
    "Uptime": ("uptime", "Server uptime in minutes"),
    "Maps": ("maps", "Maps played since startup"),
    "Map_changes": ("maps", "Maps played since startup"),
    "FPS": ("fps", "Server frames per second"),
    "Players": ("players", "Connected players"),
    "Connects": ("connects", "Connections since startup"),
    "Svms": ("svms", "ms per sim frame"),
    "+-ms": ("ms", "millisecond variance"),
    "~tick": ("tick_ms", "milliseconds per tick"),
}
STATS_HELP = dict(STATS_COLUMNS.values())

STATUS_PATTERNS = {
    "hostname": re.compile(r"hostname: (.*)"),
    "version": re.compile(r"version : (.*)"),
    "map": re.compile(r"map     : (\S+)"),
    "players": re.compile(r"players : (\d+) humans, (\d+) bots \((\d+)(?:/\d+)? max\)(?: (\((?:not )?hibernating\)))?"),
}


def parse_stats(resp: str) -> dict[str, float]:
    lines = resp.splitlines()
    for header, row in zip(lines, lines[1:]):
        if header.split()[:1] == ["CPU"]:
            return {STATS_COLUMNS[h][0]: float(v) for h, v in zip(header.split(), row.split()) if h in STATS_COLUMNS}
    raise ValueError(f"no stats table in response: {resp!r}")


def parse_status(resp: str) -> dict:
    status = {}
    for line in resp.splitlines():
        for key, pattern in STATUS_PATTERNS.items():
            if m := pattern.search(line):
                break
        else:
            continue
        if key == "players":
            humans, bots, max_players, hibernating = m.groups()
            status["players"] = int(humans)
            status["bots"] = int(bots)
            status["max_players"] = int(max_players)
            status["hibernating"] = hibernating == "(hibernating)"
        else:
            status[key] = m.group(1).strip()
    return status


class SrcdsCollector(Collector):
    def __init__(self, host: str, port: int, password: str, timeout: float = 5.0):
        self.host, self.port, self.password, self.timeout = host, port, password, timeout

    def collect(self):
        up = GaugeMetricFamily("srcds_up", "Whether the last RCON query succeeded")
        try:
            with Client(self.host, self.port, passwd=self.password, timeout=self.timeout) as client:
                stats = parse_stats(client.run("stats"))
                status = parse_status(client.run("status"))
        except Exception as e:  # Any failure, network or parsing, means the server is down to us.
            log.warning("failed to query %s:%d: %r", self.host, self.port, e)
            up.add_metric([], 0)
            yield up
            return
        up.add_metric([], 1)
        yield up

        for name, value in stats.items():
            yield GaugeMetricFamily(f"srcds_stats_{name}", STATS_HELP[name], value=value)

        if "players" in status:
            yield GaugeMetricFamily("srcds_player_count", "Connected human players", value=status["players"])
            yield GaugeMetricFamily("srcds_bot_count", "Connected bots", value=status["bots"])
            yield GaugeMetricFamily("srcds_max_players", "Player slots", value=status["max_players"])
            yield GaugeMetricFamily("srcds_active", "1 if the server is not hibernating", value=int(not status["hibernating"]))

        info = InfoMetricFamily("srcds", "Server information")
        info.add_metric([], {
            "hostname": status.get("hostname", ""),
            "version": status.get("version", ""),
            "map": status.get("map", ""),
            "active": "false" if status.get("hibernating", True) else "true",
        })
        yield info


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--listen_address", default="0.0.0.0:80", help="The address to listen on for HTTP requests.")
    parser.add_argument("--rcon_address", default="localhost:27015", help="Rcon address")
    parser.add_argument("--rcon_password", default=os.environ.get("RCON_PASSWORD", "changeme"), help="Rcon password (default: $RCON_PASSWORD)")
    args = parser.parse_args()
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")

    rcon_host, rcon_port = args.rcon_address.rsplit(":", 1)
    listen_host, listen_port = args.listen_address.rsplit(":", 1)

    registry = CollectorRegistry()
    registry.register(SrcdsCollector(rcon_host, int(rcon_port), args.rcon_password))
    _, thread = start_http_server(int(listen_port), listen_host, registry)
    log.info("listening on %s, scraping %s", args.listen_address, args.rcon_address)
    thread.join()


if __name__ == "__main__":
    main()
