package v3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/ice_cream_maker"
	"github.com/jaehoonkim/sentinel/pkg/manager/macro/slicestrings"
	"github.com/jaehoonkim/sentinel/pkg/manager/model/tenants/v3"
)

var (
	TableNameWithTenant = tableNameWithTenant()
	TenantTableName     = tenantTableName()
)

func tableNameWithTenant() func(tenant_hash string) string {
	var C = Cluster{}
	var TC = tenants.TenantClusters{}
	var T = tenants.Tenant{}

	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(C), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasC + "." + s
	})

	tables := []string{aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.id = %v.cluster_id", aliasC, aliasTC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tenantTableName() func(cluster_uuid string) string {
	var C = Cluster{}
	var TC = tenants.TenantClusters{}
	var T = tenants.Tenant{}

	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(T), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasT + "." + s
	})

	tables := []string{aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.uuid = '%%v'", aliasC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.cluster_id =  %v.id", aliasTC, aliasC),
		fmt.Sprintf("%v.id = %v.tenant_id", aliasT, aliasTC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(cluster_uuid string) string {
		return fmt.Sprintf(format, cluster_uuid)
	}
}
