class ReviewFilter {
    constructor() {
        this.initialize();
    }

    async initialize() {
        const reviewDivs = this.findReviewDivs();
        if (reviewDivs.length > 0) {
            await this.processReviews(reviewDivs);
        }
    }

    findReviewDivs() {
        return Array.from(document.querySelectorAll('div[data-test="review-details-container"]'))
            .filter(div => div.id && div.id.startsWith('empReview'));
    }

    async processReviews(reviewDivs) {
        const reviewContents = reviewDivs.map(div => {
            const ratingSpan = div.querySelector('[data-test="review-rating-label"]');
            const rating = ratingSpan ? parseFloat(ratingSpan.textContent) : 0;

            const prosSpan = div.querySelector('[data-test="review-text-PROS"]');
            const pros = prosSpan ? prosSpan.textContent.trim() : '';

            const consSpan = div.querySelector('[data-test="review-text-CONS"]');
            const cons = consSpan ? consSpan.textContent.trim() : '';

            return {
                id: div.id,
                rating: rating,
                pros: pros,
                cons: cons
            };
        });

        console.log(reviewContents)

        try {
            console.log("Worker calling");
            chrome.runtime.sendMessage( //goes to worker.js
                reviewContents,
                data => console.log(data)/*this.handleApiResponse(data, reviewDivs)*/
            );

        } catch (error) {
            console.error('Error processing reviews:', error);
        }
    }

    handleApiResponse(result, reviewDivs) {
        reviewDivs.forEach((div, index) => {
            if (result.hideReviews.includes(div.id)) {
                this.hideReview(div);
            }
        });
    }

    hideReview(reviewDiv) {
        reviewDiv.style.display = 'none';
        const flashMessage = this.createFlashMessage(reviewDiv);
        reviewDiv.parentNode.insertBefore(flashMessage, reviewDiv);
    }

    createFlashMessage(reviewDiv) {
        const flashDiv = document.createElement('div');
        flashDiv.className = 'flash-message';
        flashDiv.innerHTML = `
        <a class="show-review">Show Review</a>
        <span>This review has been hidden because it seems to be fake.</span>
      `;

        const showReviewLink = flashDiv.querySelector('.show-review');
        showReviewLink.addEventListener('click', () => this.showReview(reviewDiv, flashDiv));

        return flashDiv;
    }

    showReview(reviewDiv, flashDiv) {
        reviewDiv.style.display = 'block';
        flashDiv.innerHTML = `
        <a class="genuine-review">Seems Genuine</a>
        <a class="fake-review">Seems Fake</a>
        <span>This review has been hidden because it seems to be fake.</span>
      `;

        const genuineLink = flashDiv.querySelector('.genuine-review');
        const fakeLink = flashDiv.querySelector('.fake-review');

        genuineLink.addEventListener('click', () => this.markReviewGenuine(reviewDiv.id));
        fakeLink.addEventListener('click', () => {
            this.markReviewFake(reviewDiv.id);
            reviewDiv.style.display = 'none';
            this.resetFlashMessage(flashDiv, reviewDiv);
        });
    }

    async markReviewGenuine(reviewId) {
        try {
            await fetch(`${this.API_ENDPOINT}/genuine`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ reviewId })
            });
        } catch (error) {
            console.error('Error marking review as genuine:', error);
        }
    }

    async markReviewFake(reviewId) {
        try {
            await fetch(`${this.API_ENDPOINT}/fake`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ reviewId })
            });
        } catch (error) {
            console.error('Error marking review as fake:', error);
        }
    }

    resetFlashMessage(flashDiv, reviewDiv) {
        flashDiv.innerHTML = `
        <a class="show-review">Show Review</a>
        <span>This review has been hidden because it seems to be fake.</span>
      `;
        const showReviewLink = flashDiv.querySelector('.show-review');
        showReviewLink.addEventListener('click', () => this.showReview(reviewDiv, flashDiv));
    }
}

// Initialize the extension
new ReviewFilter();