"""End-to-end test of the exporter image against a real CS:Source server in Docker.

Run with: bazel test //:srcds_exporter_itest
"""

import os
import subprocess
import time
import unittest
import urllib.request
import uuid

from prometheus_client.parser import text_string_to_metric_families
from rcon.source import Client

SERVER_IMAGE = "ghcr.io/lanofdoom/counterstrikesource-server:latest"
EXPORTER_IMAGE = "ghcr.io/kleinpa/srcds-exporter:latest"
PASSWORD = "itest"


def docker(*args):
    return subprocess.run(["docker", *args], check=True, capture_output=True, text=True).stdout.strip()


def published_port(container, port):
    return int(docker("port", container, port).splitlines()[0].rsplit(":", 1)[1])


def eventually(check, timeout=60):
    deadline = time.monotonic() + timeout
    while True:
        try:
            return check()
        except Exception:
            if time.monotonic() > deadline:
                raise
            time.sleep(2)


class ExporterIntegrationTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        subprocess.run([os.environ["LOAD_IMAGE"]], check=True)

        suffix = uuid.uuid4().hex[:8]
        network, cls.server, exporter = f"srcds-itest-{suffix}", f"css-{suffix}", f"exporter-{suffix}"
        docker("network", "create", network)
        cls.addClassCleanup(docker, "network", "rm", network)
        docker("run", "-d", "--name", cls.server, "--network", network, "-p", "127.0.0.1::27015/tcp",
               "-e", f"RCON_PASSWORD={PASSWORD}", "-e", "CSS_HOSTNAME=exporter-itest", SERVER_IMAGE)
        cls.addClassCleanup(docker, "rm", "-f", cls.server)
        docker("run", "-d", "--name", exporter, "--network", network, "-p", "127.0.0.1::80/tcp",
               EXPORTER_IMAGE, f"--rcon_address={cls.server}:27015", f"--rcon_password={PASSWORD}")
        cls.addClassCleanup(docker, "rm", "-f", exporter)

        cls.rcon_port = published_port(cls.server, "27015/tcp")
        cls.metrics_url = f"http://127.0.0.1:{published_port(exporter, '80/tcp')}/metrics"
        eventually(lambda: cls.rcon("status"), timeout=300)  # Server installs packages on boot.

    @classmethod
    def rcon(cls, command):
        with Client("127.0.0.1", cls.rcon_port, passwd=PASSWORD, timeout=5) as client:
            return client.run(command)

    def scrape(self):
        with urllib.request.urlopen(self.metrics_url, timeout=10) as resp:
            families = text_string_to_metric_families(resp.read().decode())
        return {s.name: s for f in families for s in f.samples}

    def assertMetrics(self, expected, timeout=60):
        def check():
            samples = self.scrape()
            for name, value in expected.items():
                self.assertIn(name, samples)
                if isinstance(value, dict):
                    self.assertLessEqual(value.items(), samples[name].labels.items())
                else:
                    self.assertEqual(samples[name].value, value, name)
        eventually(check, timeout)

    def test_end_to_end(self):
        self.assertMetrics({
            "srcds_up": 1,
            "srcds_bot_count": 0,
            "srcds_player_count": 0,
            "srcds_info": {"hostname": "exporter-itest", "map": "de_dust2", "active": "true"},
        })

        self.rcon("bot_join_after_player 0")
        self.rcon("bot_quota 4")
        self.assertMetrics({"srcds_bot_count": 4, "srcds_stats_players": 4})

        self.rcon("changelevel cs_office")
        self.assertMetrics({"srcds_info": {"map": "cs_office"}, "srcds_stats_maps": 2})

        docker("stop", self.server)
        self.assertMetrics({"srcds_up": 0})


if __name__ == "__main__":
    unittest.main()
