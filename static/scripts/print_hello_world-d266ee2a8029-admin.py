#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import json

try:
    json_data = sys.argv[1]
    form_data = json.loads(sys.argv[1])
except Exception as e:
    print("waring: ", e)

# ------------- 在下面编写您的业务逻辑代码 -------------
print(json_data)
print("test ok!")