import boto3
from django.core.management.base import BaseCommand


class Command(BaseCommand):
    help = "Replicate card images from database URLs to an S3 bucket"

    def add_arguments(self, parser):
        parser.add_argument("s3_url", help="URL of target S3 bucket")
        parser.add_argument("--sets", nargs="+")

    def handle(self, *args, **options):
        client = boto3.client('s3', endpoint_url='http://localhost:9000', aws_access_key_id='minioadmin', aws_secret_access_key='minioadmin')
        resp = client.list_buckets()
        print(resp['Buckets'])
