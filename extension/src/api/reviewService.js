export class ReviewService {
    constructor(apiEndpoint) {
        this.apiEndpoint = apiEndpoint;
    }

    async checkReviews(reviews) {
        try {
            const response = await fetch(this.apiEndpoint, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ reviews })
            });
            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    async markReviewGenuine(reviewId) {
        return this.updateReviewStatus(reviewId, 'genuine');
    }

    async markReviewFake(reviewId) {
        return this.updateReviewStatus(reviewId, 'fake');
    }

    private async updateReviewStatus(reviewId, status) {
        try {
            const response = await fetch(`${this.apiEndpoint}/${status}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ reviewId })
            });
            return await response.json();
        } catch (error) {
            console.error(`Error marking review as ${status}:`, error);
            throw error;
        }
    }
}