.PHONY: run test fmt

run:
	@echo "==> Starting run on local..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' build"
	@sh -c "sh '$(CURDIR)/scripts/local.sh' run"

test:
	@echo "==> Testing..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' test"

fmt:
	@echo "==> Formatting go sources..."
	@sh -c "sh '$(CURDIR)/scripts/local.sh' fmt"

# Arguments have priority
#FOO=hoge
#export FOO
#testecho:
#	@echo $(FOO)

