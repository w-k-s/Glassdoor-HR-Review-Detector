from datetime import datetime
from pathlib import Path
from bs4 import BeautifulSoup
import csv


def extract_reviews_from_html_files():
    html_directory = Path("training-data")
    html_files = list(html_directory.glob("*.html"))

    reviews_data = []
    for html_file in html_files:
        content = html_file.read_text()
        soup = BeautifulSoup(content, "html.parser")

        reviews = soup.find_all("div", class_="review-details_reviewDetails__wSGbU")
        for review in reviews:
            reviews_data.append(
                {
                    "empReview_id": review.parent.get("id"),
                    "date": review.find(
                        "span", class_="timestamp_reviewDate__wvu2v"
                    ).text,
                    "overall_rating": review.find(
                        "span", class_="review-details_overallRating__VDxCx"
                    ).text,
                    "pros": review.find(
                        "span", attrs={"data-test": "review-text-pros"}
                    ).text,
                    "cons": review.find(
                        "span", attrs={"data-test": "review-text-cons"}
                    ).text,
                    "sounds_like_hr": 0,
                }
            )

    return reviews_data


def write_reviews_to_csv(reviews_data):
    timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
    csv_filename = f"training-data/glassdoor_reviews_{timestamp}.csv"

    with open(csv_filename, mode="w", newline="") as file:
        writer = csv.DictWriter(file, fieldnames=reviews_data[0].keys())
        writer.writeheader()
        writer.writerows(reviews_data)

    print(f"{len(reviews_data)-1} reviews written to {csv_filename}")


if __name__ == "__main__":
    reviews = extract_reviews_from_html_files()
    write_reviews_to_csv(reviews)
