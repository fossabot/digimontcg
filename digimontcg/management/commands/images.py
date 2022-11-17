from django.core.management.base import BaseCommand


class Command(BaseCommand):
    help = "Sync card images between BTCP+ and a target destination"

    def add_arguments(self, parser):
        parser.add_argument("--port", type=int, help="Local port to listen on")

    def handle(self, *args, **options):
        self.stdout.write("TODO")
