# Glassdoor HR Review Detector

My (likely to fail) attempt at building an AI that can detect if a glassdoor review was written by HR.

## Prerequisites

- Python 3.12

## Setup

1. Navigate to the project root directory. If a virtual environment does not exist, create one:

      ```shell
      python3 -m venv .venv
      ```

      Then activate the virtual environment:

      ```shell
      source .venv/bin/activate
      ```

1. Install dependencies

      ```shell
      python -m pip install --upgrade pip
      pip install -r requirements.txt

1. Start the runbook

    ```shell
    .venv/bin/jupyter notebook
    ```

## DataSets

- [HuggingFace](https://huggingface.co/datasets/lallantop/glassdoor/tree/main)
