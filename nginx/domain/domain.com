server {
         listen 7081;

         error_log logs/domain-error.log error;
         access_log logs/domain-access.log access;
         default_type text/plain;
         charset utf-8;

         #security token
         set $st "";

         #产品编号
         set $product_id "";

         #用户ID
         set_by_lua_file $user_id D:/project/Go/blitzSeckill/nginx//lua/set_common_var.lua;

         location /query {
            limit_req zone=limit_by_user nodelay;
            proxy_pass http://backend;
            #设置返回的header，并将security token放在header中
            header_filter_by_lua_block{
               ngx.header["st"] = ngx.md5(ngx.var.user_id.."1")
               ngx.header["Access-Control-Expose-Headers"] = "st"

            }
         }

        #进结算页页面（H5）
        location /settlement/prePage{
            default_type text/html;
            rewrite_by_lua_block{
                --校验活动查询的st
                local _st = ngx.md5(ngx.var.user_id.."1")
                --校验不通过时，以500状态码，返回对应错误页
                if _st ~= ngx.var.st then
                    --ngx.log(ngx.ERR, "_st value: ", _st)
                    --ngx.log(ngx.ERR, "ngx.var.st value: ", ngx.var.st)
                  ngx.log(ngx.ERR,"st is not valid!!")
                  return ngx.exit(500)
                end
                --校验通过时，再生成个新的st，用于下个接口校验
                local new_st = ngx.md5(ngx.var.user_id.."2")
                --ngx.exec执行内部跳转,浏览器URL不会发生变化
                --ngx.redirect(url,status) 其中status为301或302
                local redirect_url = "/settlement/page".."?productId="..ngx.var.product_id.."&st="..new_st
                return ngx.redirect(redirect_url,302)
            }
            error_page 500 502 503 504 /html_fail.html;
        }

        #进结算页页面（H5）
         location /settlement/page{
             default_type text/html;
             proxy_pass http://backend;
             error_page 500 502 503 504 /html_fail.html;
         }

        #结算页页面初始化渲染所需数据
        location /settlement/initData{
            access_by_lua_block{
               local _st = ngx.md5(ngx.var.user_id.."2")
               if _st ~= ngx.var.st then
                 return ngx.exit(500)
               end
            }
            proxy_pass http://backend;
            header_filter_by_lua_block{
               ngx.header["st"] = ngx.md5(ngx.var.user_id.."3")
               ngx.header["Access-Control-Expose-Headers"] = "st"
            }
            error_page 500 502 503 504 @json_fail;
        }

        #结算页提交订单
        location /settlement/submitData{
            access_by_lua_file D:/project/Go/blitzSeckill/nginx/lua/submit_access.lua;
            proxy_pass http://backend;
            error_page 500 502 503 504 @json_fail;
        }

        include D:/project/Go/blitzSeckill/nginx/domain/public.com;

}