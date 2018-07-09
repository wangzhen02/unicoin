package visor

import (
	"github.com/skycoin/skycoin/src/coin"
)

const (
	// MaxCoinSupply is the maximum supply of skycoins
	MaxCoinSupply uint64 = 1e8 // 100,000,000 million

	// DistributionAddressesTotal is the number of distribution addresses
	DistributionAddressesTotal uint64 = 100

	// DistributionAddressInitialBalance is the initial balance of each distribution address
	DistributionAddressInitialBalance uint64 = MaxCoinSupply / DistributionAddressesTotal

	// InitialUnlockedCount is the initial number of unlocked addresses
	InitialUnlockedCount uint64 = 25

	// UnlockAddressRate is the number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64 = 5

	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 = 60 * 60 * 24 * 365 // 1 year
)

func init() {
	if MaxCoinSupply%DistributionAddressesTotal != 0 {
		panic("MaxCoinSupply should be perfectly divisible by DistributionAddressesTotal")
	}
}

// GetDistributionAddresses returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 100 addresses.
func GetDistributionAddresses() []string {
	addrs := make([]string, len(distributionAddresses))
	for i := range distributionAddresses {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetUnlockedDistributionAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
func GetUnlockedDistributionAddresses() []string {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]string, InitialUnlockedCount)
	for i := range distributionAddresses[:InitialUnlockedCount] {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetLockedDistributionAddresses returns distribution addresses that are locked, i.e. they have unspendable outputs
func GetLockedDistributionAddresses() []string {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]string, DistributionAddressesTotal-InitialUnlockedCount)
	for i := range distributionAddresses[InitialUnlockedCount:] {
		addrs[i] = distributionAddresses[InitialUnlockedCount+uint64(i)]
	}
	return addrs
}

// TransactionIsLocked returns true if the transaction spends locked outputs
func TransactionIsLocked(inUxs coin.UxArray) bool {
	lockedAddrs := GetLockedDistributionAddresses()
	lockedAddrsMap := make(map[string]struct{})
	for _, a := range lockedAddrs {
		lockedAddrsMap[a] = struct{}{}
	}

	for _, o := range inUxs {
		uxAddr := o.Body.Address.String()
		if _, ok := lockedAddrsMap[uxAddr]; ok {
			return true
		}
	}

	return false
}

var distributionAddresses = [DistributionAddressesTotal]string{
	"9qZPcewnkLEfgKQ27ziYGuJkP4p3WDuQqV",
	"2MnbJ1amVWrNmmonTXH9xwwL6R9Mmg8AGub",
	"2RQwg9xFgmW3tfrkBQNf2zqubSWCarNtAdH",
	"2dtFuTTy3JBeAnQPRNomSmp3Egmts2Q5bGm",
	"uLnUhgt6kS3B7jVw6HGS4wWDwccjBpv3zN",
	"2GzFyQXFokPayyz1f4ex2qs4FLfn27CK9Hk",
	"NqMTTv2XuDPdbpJNe7idaAfH7Tw2xa9NCR",
	"MriQWGCvSyoj6C7BeW7szEaLnfurttokmx",
	"iP4qJhANFFbKrNk3rCvZYRv4nQnXJyzw6Z",
	"VX228kwbdrEkWrhhU5cjCqj3tLaoVpEVWQ",
	"wpRhQjsxHudPjecVdTW4U2TvuoA3ChwQTa",
	"2GxfnXDhjjibM7CaJAgfixc2C1WegATfEeR",
	"2LhusDfa5A73uDiquePEfXEjwePJz4f8cfB",
	"2EhTYfb2vVb4gKXRz8Vy1y2NsRS288wBw2B",
	"2jDr3L5mzBdHWwozjbkaHBZGqm7DXaVhEni",
	"26MxWZ1iqt61YvxRrHytdfsY34xz9pLj2Az",
	"2CrzTJtM19dXXJWJq4zTVRvde1kgrX2pTwW",
	"2ZE9K5qvrXBf1ACDFwMM9vr9ueLiZM8ysrX",
	"RU4i5wUsv8ANXoBJC2eDuuvejWu8rvjtCH",
	"W3uP5iEpZq1bJxNBKmZMikdL7c156LCiob",
	"2e5Dehcc8BMUeCjBXvZUvk64VZFVAzLFiF3",
	"QDMfJbNxpDqka5oKjFke6ccCUQCbHQmCcB",
	"2bxcLQYgjV14zhDw75XmumPFChafa7pGBUy",
	"25MX9z1cujYS1464WrEJkGuEv4fqcorJPsJ",
	"2NWBFzpSE17ufinPhVMy6odS7Zd8VdmNn3T",
	"2DHRn7xU4CeqC7mkscwzmgdnf3y8CL88oE7",
	"2T32LD9V4WgZuaznBHH1ediGMdmxbry8Hxi",
	"2SZm9QkwcSdacsNahsN6vdaDtbrLob4DGef",
	"9334AZxPmiKDmEzQKEgiiv1M4SSLF2Ahe2",
	"5RACpDUqnt1NmZUQ429oA74zXEbuMk4hn8",
	"YmYTKqP7zBS5VxwL5o4Sb1j8RcokZPzj8Y",
	"2ACe1p9pAfjd8mXTPq7Y4yiFtoPf3tbq9JW",
	"25MBxA8v1LmHKKCseBqya3CDo6eE3Uz25xd",
	"2KvjGdxgXMBQp9mtQF8b9QXzUTAXdykatvH",
	"2TWc2gHD9BmyMpRgRoGcQxoV4TxyHrzQKAs",
	"xnaoepH679XJyLUFUW1Wp4rqkYtV8CvTeB",
	"2GVSieJWtuBdnpk4pvQ5k7aZrTuTW8kFMVy",
	"JBMeu8B8XEDupAV7Sp5LHXCoVG1bVzrhps",
	"kpxTvDVehgxkm2ZDegpw6QGChiGCq5LuzA",
	"29WeiS2nRiZbtbwYZrwZXNL2N4oxDW6HKQF",
	"2N187MvsYKEes4Qf83nWK6h5bQNFGZvrsKG",
	"YRi6qm7EQsPuMLndtY3XXodD11No1HBN7q",
	"vLvmjCuLFbKYifHfq2UKCXDyC4kzHhU7da",
	"7H2TVj3ReJRfBXrXLQ5GcDtYWwZxbTnxSy",
	"7LBzZY8q8oHKpvcRrhmaCnTH6kzK8NiPFZ",
	"2aQJXM4VejBUh3PEe6NmQXek5ESbJ7t2324",
	"273ytDnnWco6Daubfib7oWhR76e5ucLuKfy",
	"Ykcyy2X7M6Nh1snTRsntjLyucijDDtYTRR",
	"9f4pqQt2WdHJPnvt1XyTaJjFYnH99duFHs",
	"DW9YoM7Ji1vUW9QeogYQWmtguB5qmmUCuo",
	"2BkfMoXxPWkJMCZCf8KJikrMKwYq7YBGmmt",
	"4DX7znTsjndLgfCvdNgsibwgzMJK1fwgzS",
	"rNKVDHCWfcb4BKKf9QhqnWDbVJkq3A2qW7",
	"2jYKwSLkhqkE7KnDYeKkX67VsadtxHKqJt7",
	"9pzbzY6E1ENHggztxhrtcYnWQggKzDnj45",
	"dsnAysTBYejVa4tpY8Y5xqr4iEJY6tTzRS",
	"8ad7HLaoT2Sm4EuwYUfS67YR7xrz2GUmuw",
	"FMQdm9MbP3KBEx2ioxDW2D6zNNWBAStT72",
	"2ehgPRpYN37PeRFxRFfq3GJPD2c6NCGNpMn",
	"ddVvHsQY6qgwpdfoMb3kj8MKSMfe6tq72S",
	"dTHuy4jgso9vftyUZHCj2dDktM7aVVxSVA",
	"QCwzdCQtpdTiUywkX1obzUPgcpGRLU1Jd4",
	"fyFgBFiNeJ82z3s4LUzMaRNh2PM1rfLhqW",
	"2AZvwBwVaRe88DeNvnc8gNY59Ln8jVNBcfK",
	"WEQuP35srBsTaVNjBByjLmqoyvGezK2dfK",
	"M1mA4CzD9sZUjm3Nm7Cdv1tdKWAMTiP8wz",
	"4gLz96KVoj3oMCBvCAPD1K8eoCwF2nZ1oZ",
	"2EsXtSdXxnAwt7sH4jyLDmxrLHjy1H71Gxr",
	"2YUo1CajXcUsxcDswmVwbHfH9v4gwrf6Mfa",
	"Kic1GKTGmaWgDcDGCmwozWfiZLYfqYjA2h",
	"2ZRRtDiwL7mzwoEaV5Aj3i5VjVLKKu5GgUf",
	"gxuMUu2U1N9d3mjX2xmVCR3z8CDncM14vQ",
	"i3aopHVY9XrBRb9Lm8irMWHewJLiAvA4pk",
	"2TwREkrwVNnzisLg6DjfQEQF4aq1SSPi9iQ",
	"2V7nQUMYK2jYGgDf1qoBAjUtrggM7iqBLvy",
	"9XCxPbX84vNTg6Hy8fLLY86MNpKJia5B1C",
	"1DGPNWNB6yFVZdY7r7anwayAA1uC9bH3nh",
	"o8hAWrACXzoJCAwiEhjjeTT7KyD6dAva9h",
	"CZn6BZrHUqHJS9qrYh29Uc4f8H75F3sLqv",
	"2cr4RC6BeA7TnmDt6wLkLKqQ7fmVGMEar74",
	"2etVMUjj7bwLyCEvFWYwAaVnv23QfUft43y",
	"VX3CU5R9eVC4HUiiTJtV4WnGJEcfAQsXTJ",
	"mJVar1hFgCgTDoNWRqUetTL7LkcfjmJq6Z",
	"kBiqu9ne85eU7sjN9orBV3aXeVfwwx3NPT",
	"CuXjbASdbfpUndLbuR7jamVLA75U9UVV96",
	"2YYNFFJc1ozS6GsEiGMyvxAjUimmNR7mC84",
	"tj3KF6q47Nd9FxJM362hizr5j1k7fnHVkE",
	"2fHeRKBvgxMSBUxjRV64189eewVjVDriQf2",
	"4jnEyYEs96f7vPcjHNtYjmXTP9rKPgvpjG",
	"2V3yYBM2pwPbwqRM6sGRChdAyih1tMgN7LV",
	"huXATMPHuMjA13KW5iGyg7M1afn5kRGXGe",
	"UPC891jqAKFgRyGpasQo4J4taPE3dBuYch",
	"5DWHtia4eJynLHqh2A9WKyUYn8vtAdBXmF",
	"ERv1it2LGpi2BBgEXubFMRYCFeSbxdaG2c",
	"81p4jYH7e2kT4WDu5EKiHSYRpEX9RQuquw",
	"aDGfFYpyZviW3xWF4xRabq4TXGZ13T7LFi",
	"W5jLmYiAT3XeXs9HQNsugqumGf9nKabCai",
	"2Vi6cRYvqrRiCdUiQ6DgH4rERamfLFKeRXU",
	"sEJUXSKs1D1i1q4h7pxQ1EtsiPQQXog2Xa",
	"2kdtisEC5iazCX2JEwMoHqW5QGeKYcmryXo",
}

/*
var distributionAddresses = [DistributionAddressesTotal]string{
	"R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
	"2EYM4WFHe4Dgz6kjAdUkM6Etep7ruz2ia6h",
	"25aGyzypSA3T9K6rgPUv1ouR13efNPtWP5m",
	"ix44h3cojvN6nqGcdpy62X7Rw6Ahnr3Thk",
	"AYV8KEBEAPCg8a59cHgqHMqYHP9nVgQDyW",
	"2Nu5Jv5Wp3RYGJU1EkjWFFHnebxMx1GjfkF",
	"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
	"tWZ11Nvor9parjg4FkwxNVcby59WVTw2iL",
	"m2joQiJRZnj3jN6NsoKNxaxzUTijkdRoSR",
	"8yf8PAQqU2cDj8Yzgz3LgBEyDqjvCh2xR7",
	"sgB3n11ZPUYHToju6TWMpUZTUcKvQnoFMJ",
	"2UYPbDBnHUEc67e7qD4eXtQQ6zfU2cyvAvk",
	"wybwGC9rhm8ZssBuzpy5goXrAdE31MPdsj",
	"JbM25o7kY7hqJZt3WGYu9pHZFCpA9TCR6t",
	"2efrft5Lnwjtk7F1p9d7BnPd72zko2hQWNi",
	"Syzmb3MiMoiNVpqFdQ38hWgffHg86D2J4e",
	"2g3GUmTQooLrNHaRDhKtLU8rWLz36Beow7F",
	"D3phtGr9iv6238b3zYXq6VgwrzwvfRzWZQ",
	"gpqsFSuMCZmsjPc6Rtgy1FmLx424tH86My",
	"2EUF3GPEUmfocnUc1w6YPtqXVCy3UZA4rAq",
	"TtAaxB3qGz5zEAhhiGkBY9VPV7cekhvRYS",
	"2fM5gVpi7XaiMPm4i29zddTNkmrKe6TzhVZ",
	"ix3NDKgxfYYANKAb5kbmwBYXPrkAsha7uG",
	"2RkPshpFFrkuaP98GprLtgHFTGvPY5e6wCK",
	"Ak1qCDNudRxZVvcW6YDAdD9jpYNNStAVqm",
	"2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf",
	"KPfqM6S96WtRLMuSy4XLfVwymVqivdcDoM",
	"5B98bU1nsedGJBdRD5wLtq7Z8t8ZXio8u5",
	"2iZWk5tmBynWxj2PpAFyiZzEws9qSnG3a6n",
	"XUGdPaVnMh7jtzPe3zkrf9FKh5nztFnQU5",
	"hSNgHgewJme8uaHrEuKubHYtYSDckD6hpf",
	"2DeK765jLgnMweYrMp1NaYHfzxumfR1PaQN",
	"orrAssY5V2HuQAbW9K6WktFrGieq2m23pr",
	"4Ebf4PkG9QEnQTm4MVvaZvJV6Y9av3jhgb",
	"7Uf5xJ3GkiEKaLxC2WmJ1t6SeekJeBdJfu",
	"oz4ytDKbCqpgjW3LPc52pW2CaK2gxCcWmL",
	"2ex5Z7TufQ5Z8xv5mXe53fSQRfUr35SSo7Q",
	"WV2ap7ZubTxeDdmEZ1Xo7ufGMkekLWikJu",
	"ckCTV4r1pNuz6j2VBRHhaJN9HsCLY7muLV",
	"MXJx96ZJVSjktgeYZpVK8vn1H3xWP8ooq5",
	"wyQVmno9aBJZmQ99nDSLoYWwp7YDJCWsrH",
	"2cc9wKxCsFNRkoAQDAoHke3ZoyL1mSV14cj",
	"29k9g3F5AYfVaa1joE1PpZjBED6hQXes8Mm",
	"2XPLzz4ZLf1A9ykyTCjW5gEmVjnWa8CuatH",
	"iH7DqqojTgUn2JxmY9hgFp165Nk7wKfan9",
	"RJzzwUs3c9C8Y7NFYzNfFoqiUKeBhBfPki",
	"2W2cGyiCRM4nwmmiGPgMuGaPGeBzEm7VZPn",
	"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
	"tBaeg9zE2sgmw5ZQENaPPYd6jfwpVpGTzS",
	"2hdTw5Hk3rsgpZjvk8TyKcCZoRVXU5QVrUt",
	"A1QU6jKq8YgTP79M8fwZNHUZc7hConFKmy",
	"q9RkXoty3X1fuaypDDRUi78rWgJWYJMmpJ",
	"2Xvm6is5cAPA85xnSYXDuAqiRyoXiky5RaD",
	"4CW2CPJEzxhn2PS4JoSLoWGL5QQ7dL2eji",
	"24EG6uTzL7DHNzcwsygYGRR1nfu5kco7AZ1",
	"KghGnWw5fppTrqHSERXZf61yf7GkuQdCnV",
	"2WojewRA3LbpyXTP9ANy8CZqJMgmyNm3MDr",
	"2BsMfywmGV3M2CoDA112Rs7ZBkiMHfy9X11",
	"kK1Q4gPyYfVVMzQtAPRzL8qXMqJ67Y7tKs",
	"28J4mx8xfUtM92DbQ6i2Jmqw5J7dNivfroN",
	"gQvgyG1djgtftoCVrSZmsRxr7okD4LheKw",
	"3iFGBKapAWWzbiGFSr5ScbhrEPm6Esyvia",
	"NFW2akQH2vu7AqkQXxFz2P5vkXTWkSqrSm",
	"2MQJjLnWRp9eHh6MpCwpiUeshhtmri12mci",
	"2QjRQUMyL6iodtHP9zKmxCNYZ7k3jxtk49C",
	"USdfKy7B6oFNoauHWMmoCA7ND9rHqYw2Mf",
	"cA49et9WtptYHf6wA1F8qqVgH3kS5jJ9vK",
	"qaJT9TjcMi46sTKcgwRQU8o5Lw2Ea1gC4N",
	"22pyn5RyhqtTQu4obYjuWYRNNw4i54L8xVr",
	"22dkmukC6iH4FFLBmHne6modJZZQ3MC9BAT",
	"z6CJZfYLvmd41GRVE8HASjRcy5hqbpHZvE",
	"GEBWJ2KpRQDBTCCtvnaAJV2cYurgXS8pta",
	"oS8fbEm82cprmAeineBeDkaKd7QownDZQh",
	"rQpAs1LVQdphyj9ipEAuukAoj9kNpSP8cM",
	"6NSJKsPxmqipGAfFFhUKbkopjrvEESTX3j",
	"cuC68ycVXmD2EBzYFNYQ6akhKGrh3FGjSf",
	"bw4wtYU8toepomrhWP2p8UFYfHBbvEV425",
	"HvgNmDz5jD39Gwmi9VfDY1iYMhZUpZ8GKz",
	"SbApuZAYquWP3Q6iD51BcMBQjuApYEkRVf",
	"2Ugii5yxJgLzC59jV1vF8GK7UBZdvxwobeJ",
	"21N2iJ1qnQRiJWcEqNRxXwfNp8QcmiyhtPy",
	"9TC4RGs6AtFUsbcVWnSoCdoCpSfM66ALAc",
	"oQzn55UWG4iMcY9bTNb27aTnRdfiGHAwbD",
	"2GCdwsRpQhcf8SQcynFrMVDM26Bbj6sgv9M",
	"2NRFe7REtSmaM2qAgZeG45hC8EtVGV2QjeB",
	"25RGnhN7VojHUTvQBJA9nBT5y1qTQGULMzR",
	"26uCBDfF8E2PJU2Dzz2ysgKwv9m4BhodTz9",
	"Wkvima5cF7DDFdmJQqcdq8Syaq9DuAJJRD",
	"286hSoJYxvENFSHwG51ZbmKaochLJyq4ERQ",
	"FEGxF3HPoM2HCWHn82tyeh9o7vEQq5ySGE",
	"h38DxNxGhWGTq9p5tJnN5r4Fwnn85Krrb6",
	"2c1UU8J6Y3kL4cmQh21Tj8wkzidCiZxwdwd",
	"2bJ32KuGmjmwKyAtzWdLFpXNM6t83CCPLq5",
	"2fi8oLC9zfVVGnzzQtu3Y3rffS65Hiz6QHo",
	"TKD93RxFr2Am44TntLiJQus4qcEwTtvEEQ",
	"zMDywYdGEDtTSvWnCyc3qsYHWwj9ogws74",
	"25NbotTka7TwtbXUpSCQD8RMgHKspyDubXJ",
	"2ayCELBERubQWH5QxUr3cTxrYpidvUAzsSw",
	"RMTCwLiYDKEAiJu5ekHL1NQ8UKHi5ozCPg",
	"ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm",
}
*/
