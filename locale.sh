# update locale
# 依赖 github.com/issue9/web

web locale -l=und -f=yaml ./
web update-locale -src=./locales/und.yaml -dest=./locales/zh.yaml
