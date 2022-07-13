
gen-mock:
	mockery --dir internal/services --all --output ./mock/services --outpkg services
