import csv
import io


def load_csv_as_dictionary(csv_filename):
    with open(csv_filename, mode="r", newline="") as file:
        reader = csv.DictReader(file)
        data = [row for row in reader]
    return data


def array_of_dictionaries_to_csv_string(items):
    with io.StringIO() as output:
        writer = csv.DictWriter(output, fieldnames=items[0].keys())
        writer.writeheader()
        writer.writerows(items)

        result = output.getvalue()
    return result
