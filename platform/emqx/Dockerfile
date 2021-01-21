FROM emqx/emqx:4.2.6
COPY ./emqx.conf /opt/emqx/etc/emqx.conf
COPY ./acl.conf /opt/emqx/etc/acl.conf
COPY ./emqx_auth_pgsql.conf /opt/emqx/etc/plugins/emqx_auth_pgsql.conf
COPY loaded_plugins /opt/emqx/data/loaded_plugins