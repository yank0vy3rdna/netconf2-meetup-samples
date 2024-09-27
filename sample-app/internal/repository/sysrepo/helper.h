#ifndef _GO_SYSREPO_LIB_H
#define _GO_SYSREPO_LIB_H

#include <sysrepo.h>

#ifdef __cplusplus
extern "C" {
#endif
  int _sr_module_change_cb(sr_session_ctx_t *session, uint32_t sub_id, const char *module_name, const char *xpath,
                           sr_event_t event, uint32_t request_id, void *private_data);

#ifdef __cplusplus
}
#endif

#endif // _GO_SYSREPO_LIB_H
