package prom_query_tool

import (
	"context"
	"fmt"
	"github.com/prometheus/common/model"
	"testing"
	"time"
)

type testWorkloadLabel struct {
	WorkloadName  string
	WorkloadType  string
	Namespace     string
	ContainerName string
}

func (t *testWorkloadLabel) String() string {
	return fmt.Sprintf("[%s]下的工作负载类型为[%s]名为[%s]的容器[%s]", t.Namespace, t.WorkloadType, t.WorkloadName, t.ContainerName)
}

const promAddress = "http://prometheus-test.kubeease.cn"

func TestPromVector(t *testing.T) {

	pq, err := NewPromQuery(promAddress)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	value, _, err := pq.Query().PromQl("max(quantile_over_time(%s,container_memory_working_set_bytes{image!=\"\", container!=\"POD\"}[%dh:%dm]) * on (namespace,pod) group_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s} ) by (namespace,workload,workload_type,container_name)", "0.99", 168, 10, "").
		Time(time.Now()).DoQuery(context.Background())
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resultMap, err := NewVector[testWorkloadLabel, string](value).SetLabelCovertFunc(func(metric model.Metric) testWorkloadLabel {
		return testWorkloadLabel{
			Namespace:     string(metric["namespace"]),
			WorkloadName:  string(metric["workload"]),
			WorkloadType:  string(metric["workload_type"]),
			ContainerName: string(metric["container_name"]),
		}
	}).SetValueConvertFunc(func(sample model.SampleValue) string {
		return fmt.Sprintf("%dKi", int64((float64(sample)/1024)+0.5))
	}).ToMap()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for k, v := range resultMap {
		fmt.Printf("%s当前使用最大内存为：%s\n", k.String(), v)
	}
}

func TestPromMatrix(t *testing.T) {
	pq, err := NewPromQuery(promAddress)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	value, _, err := pq.Query().PromQl("1- avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) by (node)").
		Range(time.Now().Add(-7*24*time.Hour), time.Now(), time.Hour*4).DoQueryRange(context.Background())
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resultMap, err := NewMatrix[string, string](value).SetLabelCovertFunc(func(metric model.Metric) string { return string(metric["node"]) }).SetValueConvertFunc(func(sample model.SampleValue) string {
		return fmt.Sprintf("%.2f%%", float64(sample)*100)
	}).ToMap()
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for k, v := range resultMap {
		fmt.Printf("节点[%s]7天的CPU使用率：%s\n", k, v)
	}
}

func TestPromScalar(t *testing.T) {
	pq, err := NewPromQuery(promAddress)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	value, _, err := pq.Query().PromQl("scalar(sum(node_cpu_seconds_total))").
		Time(time.Now()).DoQuery(context.Background())
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(NewScalar(value).GetValue())

}
