"""Prometheus exporter for Source Dedicated Server (srcds) via the public A2S query protocol."""

import argparse
import logging

import a2s
from prometheus_client import start_http_server
from prometheus_client.core import CollectorRegistry, GaugeMetricFamily, InfoMetricFamily
from prometheus_client.registry import Collector

log = logging.getLogger("srcds_exporter")


class SrcdsCollector(Collector):
    def __init__(self, host: str, port: int, timeout: float = 3.0):
        self.address, self.timeout = (host, port), timeout

    def collect(self):
        up = GaugeMetricFamily("srcds_up", "Whether the last A2S query succeeded")
        try:
            info = a2s.info(self.address, timeout=self.timeout)
        except Exception as e:  # Any failure, network or parsing, means the server is down to us.
            log.warning("failed to query %s:%d: %r", *self.address, e)
            up.add_metric([], 0)
            yield up
            return
        up.add_metric([], 1)
        yield up

        # A2S counts bots as players.
        yield GaugeMetricFamily("srcds_player_count", "Connected human players", value=info.player_count - info.bot_count)
        yield GaugeMetricFamily("srcds_bot_count", "Connected bots", value=info.bot_count)
        yield GaugeMetricFamily("srcds_max_players", "Player slots", value=info.max_players)
        yield GaugeMetricFamily("srcds_ping_seconds", "A2S query round-trip time", value=info.ping)

        server = InfoMetricFamily("srcds", "Server information")
        server.add_metric([], {"hostname": info.server_name, "version": info.version, "map": info.map_name, "game": info.game})
        yield server


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--listen_address", default="0.0.0.0:80", help="The address to listen on for HTTP requests.")
    parser.add_argument("--server_address", default="localhost:27015", help="Game server address (UDP)")
    args = parser.parse_args()
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")

    server_host, server_port = args.server_address.rsplit(":", 1)
    listen_host, listen_port = args.listen_address.rsplit(":", 1)

    registry = CollectorRegistry()
    registry.register(SrcdsCollector(server_host, int(server_port)))
    _, thread = start_http_server(int(listen_port), listen_host, registry)
    log.info("listening on %s, querying %s", args.listen_address, args.server_address)
    thread.join()


if __name__ == "__main__":
    main()
