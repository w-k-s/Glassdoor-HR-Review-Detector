import { ReviewService } from '../api/reviewService';
import { FlashMessage } from './flashMessage';

export class ReviewFilter {
    constructor() {
        this.reviewService = new ReviewService(process.env.API_ENDPOINT);
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
            .filter(div => div.id?.startsWith('empReview'));
    }

    async processReviews(reviewDivs) {
        const reviewContents = reviewDivs.map(div => ({
            id: div.id,
            content: div.textContent
        }));

        try {
            const result = await this.reviewService.checkReviews(reviewContents);
            this.handleApiResponse(result, reviewDivs);
        } catch (error) {
            console.error('Error processing reviews:', error);
        }
    }

    handleApiResponse(result, reviewDivs) {
        reviewDivs.forEach(div => {
            if (result.hideReviews.includes(div.id)) {
                this.hideReview(div);
            }
        });
    }

    hideReview(reviewDiv) {
        reviewDiv.style.display = 'none';
        const flashMessage = new FlashMessage(
            () => this.showReview(reviewDiv),
            () => this.markReviewGenuine(reviewDiv.id),
            () => this.markReviewFake(reviewDiv)
        ).create();
        reviewDiv.parentNode.insertBefore(flashMessage, reviewDiv);
    }

    showReview(reviewDiv) {
        const flashDiv = reviewDiv.previousSibling;
        reviewDiv.style.display = 'block';
        new FlashMessage(
            null,
            () => this.markReviewGenuine(reviewDiv.id),
            () => this.markReviewFake(reviewDiv)
        ).showReviewMode(flashDiv);
    }

    async markReviewGenuine(reviewId) {
        await this.reviewService.markReviewGenuine(reviewId);
    }

    async markReviewFake(reviewDiv) {
        const flashDiv = reviewDiv.previousSibling;
        await this.reviewService.markReviewFake(reviewDiv.id);
        reviewDiv.style.display = 'none';
        new FlashMessage(
            () => this.showReview(reviewDiv),
            null,
            null
        ).reset(flashDiv);
    }
}

