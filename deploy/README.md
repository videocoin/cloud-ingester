```
helm upgrade -i --wait --set image.tag="v1.0.0-pe-1-06d8149" --set hookd.image.tag="v1.0.0-pe-1-06d8149" -n console ingester ./helm
```