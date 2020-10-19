.PHONY: run test fmt

run:
	@echo "==> Starting run on local..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' build"
	@sh -c "sh '$(CURDIR)/scripts/local.sh' run"

dbrun:
	@echo "==> Starting DB on local..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' docker:run"

dbstop:
	@echo "==> Starting DBon local..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' docker:stop"

test:
	@echo "==> Testing..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' test"

testnocache:
	@echo "==> Testing..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' test nocache"

fmt:
	@echo "==> Formatting go sources..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' fmt"

# Arguments have priority
#FOO=hoge
#export FOO
#testecho:
#	@echo $(FOO)

