#include "helper.h"

extern int go_sr_module_change_cb(sr_session_ctx_t *session, uint32_t sub_id, const char *module_name, const char *xpath,
                                                                       sr_event_t event, uint32_t request_id, void *private_data);

int _sr_module_change_cb(sr_session_ctx_t *session, uint32_t sub_id, const char *module_name, const char *xpath,
                                     sr_event_t event, uint32_t request_id, void *private_data) {
    return go_sr_module_change_cb(session, sub_id, module_name, xpath, event, request_id, private_data);
}
