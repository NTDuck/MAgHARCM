package agents_test

import "testing"

func TestTryInOrderEmptyProfileRunner(t *testing.T)        { TestTryInOrderEmptyProfile(t) }
func TestTryInOrderPilotProfileRunner(t *testing.T)        { TestTryInOrderPilotProfile(t) }
func TestTryInOrderFrozenLegacyProfileRunner(t *testing.T) { TestTryInOrderFrozenLegacyProfile(t) }
func TestTryInOrderBigBangProfileRunner(t *testing.T)      { TestTryInOrderBigBangProfile(t) }
func TestTryInOrderParallelCutoverProfileRunner(t *testing.T) {
	TestTryInOrderParallelCutoverProfile(t)
}
func TestSelectAndTryStrategiesRunner(t *testing.T) { TestSelectAndTryStrategies(t) }
