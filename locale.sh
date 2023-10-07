# update locale
# 依赖 gitub.com/issue9/web

web locale -l=en-US -f=yaml ./
web update-locale -src=./locales/en-US.yaml -dest=./locales/zh-CN.yaml