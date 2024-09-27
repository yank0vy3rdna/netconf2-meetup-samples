package sysrepo

/*
#include "helper.h"
*/
import "C"

type ConnFlag C.sr_conn_flag_t

func (v ConnFlag) C() C.sr_conn_flag_t {
	return (C.sr_conn_flag_t)(v)
}

const (
	CONN_DEFAULT       ConnFlag = C.SR_CONN_DEFAULT
	CONN_CACHE_RUNNING ConnFlag = C.SR_CONN_CACHE_RUNNING
)

type SubscrFlag C.sr_subscr_flag_t

func (v SubscrFlag) C() C.sr_subscr_flag_t {
	return (C.sr_subscr_flag_t)(v)
}

const (
	SUBSCR_DEFAULT    SubscrFlag = C.SR_SUBSCR_DEFAULT
	SUBSCR_NO_THREAD  SubscrFlag = C.SR_SUBSCR_NO_THREAD
	SUBSCR_PASSIVE    SubscrFlag = C.SR_SUBSCR_PASSIVE
	SUBSCR_DONE_ONLY  SubscrFlag = C.SR_SUBSCR_DONE_ONLY
	SUBSCR_ENABLED    SubscrFlag = C.SR_SUBSCR_ENABLED
	SUBSCR_UPDATE     SubscrFlag = C.SR_SUBSCR_UPDATE
	SUBSCR_OPER_MERGE SubscrFlag = C.SR_SUBSCR_OPER_MERGE
)

type NotifyEvent C.sr_event_t

func (v NotifyEvent) C() C.sr_event_t {
	return (C.sr_event_t)(v)
}

type LydFormat C.LYD_FORMAT

const (
	LYD_UNKNOWN LydFormat = C.LYD_UNKNOWN
	LYD_XML     LydFormat = C.LYD_XML
	LYD_JSON    LydFormat = C.LYD_JSON
	LYD_LYB     LydFormat = C.LYD_LYB
)

func (f LydFormat) C() C.LYD_FORMAT {
	return C.LYD_FORMAT(f)
}

type Option C.int

const (
	LYD_PRINT_WITHSIBLINGS  Option = C.LYD_PRINT_WITHSIBLINGS
	LYD_PRINT_SHRINK        Option = C.LYD_PRINT_SHRINK
	LYD_PRINT_KEEPEMPTYCONT Option = C.LYD_PRINT_KEEPEMPTYCONT
	LYD_PRINT_WD_MASK       Option = C.LYD_PRINT_WD_MASK
	LYD_PRINT_WD_EXPLICIT   Option = C.LYD_PRINT_WD_EXPLICIT
	LYD_PRINT_WD_TRIM       Option = C.LYD_PRINT_WD_TRIM
	LYD_PRINT_WD_ALL        Option = C.LYD_PRINT_WD_ALL
	LYD_PRINT_WD_ALL_TAG    Option = C.LYD_PRINT_WD_ALL_TAG
	LYD_PRINT_WD_IMPL_TAG   Option = C.LYD_PRINT_WD_IMPL_TAG
)

func (o Option) C() C.uint32_t {
	return C.uint32_t(o)
}
