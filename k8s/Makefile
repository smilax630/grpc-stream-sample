APPS := stream

INSTALL_TARGETS := $(addprefix install, $(APPS))
UPGRADE_TARGETS := $(addprefix upgrade-, $(APPS))
UNINSTALL_TARGETS := $(addprefix uninstall-, $(APPS))

.PHONY: $(DEPLOY_DEV_TARGETS) 

$(INSTALL_TARGETS):  install-%:
	helm install $* ./$* -f $*/env.yaml

$(UPGRADE_TARGETS):  upgrade-%: 
	helm upgrade $* ./$* -f $*/env.yaml

$(UNINSTALL_TARGETS): uninstall-%:
	helm uninstall $*
