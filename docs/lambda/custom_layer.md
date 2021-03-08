# Add custom modules to AWS Lambda
```
mkdir staging
cd staging
cp -r ../src python
zip -r varity-scraping.zip python
aws lambda publish-layer-version --layer-name varity-scraping --zip-file fileb://varity-scraping.zip --compatible-runtimes python3.6 python3.7 python3.8 --description "Custom python module for varity scraping"
```