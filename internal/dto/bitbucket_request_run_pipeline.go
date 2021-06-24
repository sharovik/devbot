package dto

//PipelineTargetSelector pipeline selector
type PipelineTargetSelector struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

//Variable the object for pipeline variables
type Variable struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Secured bool   `json:"secured,omitempty"`
}

//PipelineTarget the target for the pipeline to run
type PipelineTarget struct {
	RefName   string                 `json:"ref_name"`
	RefType   string                 `json:"ref_type"`
	Selector  PipelineTargetSelector `json:"selector"`
	Type      string                 `json:"type"`
}

//BitBucketRequestRunPipeline the request for the pipeline run endpoint
type BitBucketRequestRunPipeline struct {
	Target PipelineTarget `json:"target"`
	Variables []Variable `json:"variables"`
}
