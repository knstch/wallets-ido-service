package wallets_test

import (
	"context"
)

func (s *WalletsServiceTestSuite) TestAddWallets_HappyPath() {
	t := s.Require()
	pubkey1, _ := mustGenerateSolanaKeypair(t)
	pubkey2, _ := mustGenerateSolanaKeypair(t)

	err := s.svc.AddWallets(context.Background(), 1, []string{pubkey1, pubkey2})
	t.NoError(err)

	// Verify wallets were created
	wallets, err := s.svc.GetWallets(context.Background(), 1)
	t.NoError(err)
	t.Len(wallets, 2)
	// Check that both pubkeys are in the wallets
	pubkeys := make([]string, 0, len(wallets))
	for _, w := range wallets {
		pubkeys = append(pubkeys, w.Pubkey)
	}
	t.Contains(pubkeys, pubkey1)
	t.Contains(pubkeys, pubkey2)
}

func (s *WalletsServiceTestSuite) TestAddWallets_SkipExistingForSameUser() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)

	// Add wallet first time
	err := s.svc.AddWallets(context.Background(), 1, []string{pubkey})
	t.NoError(err)

	// Add same wallet again for same user - should be skipped
	err = s.svc.AddWallets(context.Background(), 1, []string{pubkey})
	t.NoError(err)

	// Verify wallet still exists
	wallets, err := s.svc.GetWallets(context.Background(), 1)
	t.NoError(err)
	t.Len(wallets, 1)
	t.Equal(pubkey, wallets[0].Pubkey)
}

func (s *WalletsServiceTestSuite) TestAddWallets_AllowsSamePubkeyForDifferentUsers() {
	t := s.Require()
	pubkey, _ := mustGenerateSolanaKeypair(t)

	// Add wallet for user 1
	err := s.svc.AddWallets(context.Background(), 1, []string{pubkey})
	t.NoError(err)

	// Add same pubkey for user 2 - should be allowed
	err = s.svc.AddWallets(context.Background(), 2, []string{pubkey})
	t.NoError(err)

	// Verify both users have the wallet
	wallets1, err := s.svc.GetWallets(context.Background(), 1)
	t.NoError(err)
	t.Len(wallets1, 1)
	t.Equal(pubkey, wallets1[0].Pubkey)

	wallets2, err := s.svc.GetWallets(context.Background(), 2)
	t.NoError(err)
	t.Len(wallets2, 1)
	t.Equal(pubkey, wallets2[0].Pubkey)
}

func (s *WalletsServiceTestSuite) TestAddWallets_EmptyList() {
	t := s.Require()

	err := s.svc.AddWallets(context.Background(), 1, []string{})
	t.NoError(err)
}

func (s *WalletsServiceTestSuite) TestAddWallets_MixedExistingAndNew() {
	t := s.Require()
	pubkey1, _ := mustGenerateSolanaKeypair(t)
	pubkey2, _ := mustGenerateSolanaKeypair(t)

	// Add first wallet
	err := s.svc.AddWallets(context.Background(), 1, []string{pubkey1})
	t.NoError(err)

	// Add both existing and new wallet
	err = s.svc.AddWallets(context.Background(), 1, []string{pubkey1, pubkey2})
	t.NoError(err)

	// Verify both wallets exist
	wallets, err := s.svc.GetWallets(context.Background(), 1)
	t.NoError(err)
	t.Len(wallets, 2)
	// Check that both pubkeys are in the wallets
	pubkeys := make([]string, 0, len(wallets))
	for _, w := range wallets {
		pubkeys = append(pubkeys, w.Pubkey)
	}
	t.Contains(pubkeys, pubkey1)
	t.Contains(pubkeys, pubkey2)
}
