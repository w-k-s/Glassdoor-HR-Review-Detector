FROM tensorflow/serving

COPY  glassdoor-hr-review-detector/models /models/glassdoor_hr_review_detector

ENV MODEL_NAME glassdoor_hr_review_detector

ENTRYPOINT ["/usr/bin/tf_serving_entrypoint.sh"]