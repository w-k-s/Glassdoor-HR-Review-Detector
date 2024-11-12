const API_ENDPOINT = 'http://localhost:8000/api/reviews/genuity-check';

// https://stackoverflow.com/a/56502776
chrome.runtime.onMessage.addListener(
    function (data, sender, onSuccess) {
        fetch(API_ENDPOINT, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'apikey': 'client'
            },
            body: JSON.stringify({ reviews: data })
        })
            .then(response => response.json())
            .then(responseText => onSuccess(responseText))

        return true;  // Will respond asynchronously.
    }
);