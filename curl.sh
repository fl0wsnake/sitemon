curl 'https://metro.zakaz.ua/uk/products/ukrayina--04820254610291/' \
  -H 'cookie: experiment_favorite_button=1; __zlcmid=1Nen0O99eTrsyd9; storeId=48215637; deliveryType=plan' \
  | pup '.BigProductCardTopInfo__priceInfo'

curl 'https://auchan.zakaz.ua/uk/products/ukrayina--04820254610291/' \
  -H 'cookie: experiment_favorite_button=1; __zlcmid=1Nen0O99eTrsyd9; storeId=48246409; deliveryType=plan' \
  | pup '.BigProductCardTopInfo__priceInfo'
