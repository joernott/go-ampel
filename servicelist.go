package main

type ServiceList []*Service

// Finds a service by its name
func (list ServiceList) Find(serviceName string) *Service {
	var svc *Service
	for _, svc = range list {
		if svc.Name == serviceName {
			return svc
		}
	}
	return &Service{Name: ""}
}

// Appends a service entry to the list
func (list ServiceList) Append(NewService *Service) ServiceList {
	list = append(list, NewService)
	return list
}

// Lists services by status
func (ServiceList) StatusList(status bool) []string {
	var output []string

	for _, svc := range Services {
		if svc.Status == status {
			output = append(output, svc.Name)
		}
	}
	return output
}
