
terraform {

}

provider "google" {
    project = "sacred-lane-281804"
    region  = "us-central1"
}

// code storage bucket, etc
resource "google_storage_bucket" "deployment-bucket" {
    name = "moz-gcp-functions"
    force_destroy = true
}

data "archive_file" "pubsub-trigger" {
    type = "zip"
    output_path = "${path.module}/files/retweeter.zip"

    source {
        content = "${file("${path.module}/../functions/retweeter.go")}"
        filename = "main.go"
    }

    source {
        content = "${file("${path.module}/../functions/go.mod")}"
        filename = "go.mod"
    }

    source {
        content = "${file("${path.module}/../functions/go.sum")}"
        filename = "go.sum"
    }
}

resource "google_storage_bucket_object" "archive" {
    bucket = "${google_storage_bucket.deployment-bucket.name}"
    name = "${format("retweeter-%s", data.archive_file.pubsub-trigger.output_md5)}"
    source = "${path.module}/files/retweeter.zip"
    depends_on = [data.archive_file.pubsub-trigger]
}

// this topic is really here because we can't schedule a lambda but we can schedule to push
// to pub/sub
resource "google_pubsub_topic" "topic" {
  name = "retweeter-topic"
}

resource "google_cloudfunctions_function" "retweeter" {
    name = "moz-retweeter"
    runtime = "go113"

    available_memory_mb = 128
    source_archive_bucket = google_storage_bucket.deployment-bucket.name
    source_archive_object = google_storage_bucket_object.archive.name
    entry_point = "Retweeter"

    event_trigger {
        event_type = "google.pubsub.topic.publish"
        resource = "${google_pubsub_topic.topic.name}"
    }

    labels = {
        last_updated_at = "hi1"
    }

    depends_on = [ google_storage_bucket_object.archive ]
}

// schedule to push a message to pubsub every hour
resource "google_cloud_scheduler_job" "job" {
    name = "retweeter-job"
    schedule    = "1 */1 * * *"
    time_zone = "CST"

    pubsub_target {
        topic_name =   google_pubsub_topic.topic.id
        data = base64encode("a")
    }

    depends_on = [google_cloudfunctions_function.retweeter]
}
