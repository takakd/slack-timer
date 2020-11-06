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
	@sh -c "sh '$(CURDIR)/scripts/local.sh' test -v -cover -tags=\'test local\' -count 1 ./..."

test_light:
	@echo "==> Testing..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' test -tags=\'test local\' -count 1 ./..."

fmt:
	@echo "==> Formatting go sources..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' fmt"

# Arguments have priority
#FOO=hoge
#export FOO
#testecho:
#	@echo $(FOO)

