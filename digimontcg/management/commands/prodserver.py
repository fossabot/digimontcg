from django.core.management.base import BaseCommand
import waitress

from project.wsgi import application


class Command(BaseCommand):
    help = "Run the app with a production-ready WSGI server"

    def add_arguments(self, parser):
        parser.add_argument("--port", type=int, help="Local port to listen on")

    def handle(self, *args, **options):
        port = options["port"] or "5000"
        port = int(port)
        self.stdout.write(f"Listening on port {port}...")
        waitress.serve(application, host="127.0.0.1", port=port)
