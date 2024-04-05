from diagrams import Diagram
from diagrams.onprem.database import Clickhouse
from diagrams.onprem.queue import RabbitMQ
from diagrams.onprem.monitoring import Grafana
from diagrams.programming.language import Go
from diagrams.programming.language import C

with Diagram("Scheme of collecting sensors data", show=False):
    C("Arduino sensors") >> RabbitMQ("RabbitMQ") >> Go("Consumer") >> Clickhouse("Clickhouse") << Grafana("Grafana")