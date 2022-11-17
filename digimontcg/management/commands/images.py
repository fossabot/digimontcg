from django.core.management.base import BaseCommand, CommandError
from minio import Minio


class Command(BaseCommand):
    help = "Replicate card images from database URLs to an S3 bucket"

    def add_arguments(self, parser):
        parser.add_argument("s3_endpoint", help="URL of S3 endpoint")
        parser.add_argument("s3_bucket", help="Name of S3 bucket")
        parser.add_argument("--sets", nargs="+")

    def handle(self, *args, **options):
        endpoint = options["s3_endpoint"]
        bucket = options["s3_bucket"]

        client = Minio(
            endpoint,
            access_key="minioadmin",
            secret_key="minioadmin",
            secure=False,
        )

        found = client.bucket_exists(bucket)
        if not found:
            client.make_bucket(bucket)
        else:
            print(f"Bucket '{bucket}' already exists")
