from django.core.management.base import BaseCommand


class Command(BaseCommand):
    help = "Replicate card images from database URLs to an S3 bucket"

    def add_arguments(self, parser):
        parser.add_argument("s3_url", type=int, help="URL of target S3 bucket")
        parser.add_argument("--sets", nargs="+")

    def handle(self, *args, **options):
        self.stdout.write("TODO")
