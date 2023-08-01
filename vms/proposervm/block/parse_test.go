// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"crypto"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
)

func TestParse(t *testing.T) {
	require := require.New(t)

	parentID := ids.ID{1}
	timestamp := time.Unix(123, 0)
	pChainHeight := uint64(2)
	innerBlockBytes := []byte{3}
	chainID := ids.ID{4}

	tlsCert, err := staking.NewTLSCert()
	require.NoError(err)

	cert := tlsCert.Leaf
	key := tlsCert.PrivateKey.(crypto.Signer)

	builtBlock, err := Build(
		parentID,
		timestamp,
		pChainHeight,
		cert,
		innerBlockBytes,
		chainID,
		key,
	)
	require.NoError(err)

	builtBlockBytes := builtBlock.Bytes()

	parsedBlockIntf, err := Parse(builtBlockBytes)
	require.NoError(err)

	parsedBlock, ok := parsedBlockIntf.(SignedBlock)
	require.True(ok)

	equal(require, chainID, builtBlock, parsedBlock)
}

func TestParseDuplicateExtension(t *testing.T) {
	require := require.New(t)

	blockHex := "0000000000000100000000000000000000000000000000000000000000000000000000000000000000000000007b0000000000000002000004bd308204b9308202a1a003020102020100300d06092a864886f70d01010b050030003020170d3939313233313030303030305a180f32313232303830333233323835335a300030820222300d06092a864886f70d01010105000382020f003082020a0282020100c2b2de1c16924d9b9254a0d5b80a4bc5f9beaa4f4f40a0e4efb69eb9b55d7d37f8c82328c237d7c5b451f5427b487284fa3f365f9caa53c7fcfef8d7a461d743bd7d88129f2da62b877ebe9d6feabf1bd12923e6c12321382c782fc3bb6b6cb4986a937a1edc3814f4e621e1a62053deea8c7649e43edd97ab6b56315b00d9ab5026bb9c31fb042dc574ba83c54e720e0120fcba2e8a66b77839be3ece0d4a6383ef3f76aac952b49a15b65e18674cd1340c32cecbcbaf80ae45be001366cb56836575fb0ab51ea44bf7278817e99b6b180fdd110a49831a132968489822c56692161bbd372cf89d9b8ee5a734cff15303b3a960ee78d79e76662a701941d9ec084429f26707f767e9b1d43241c0e4f96655d95c1f4f4aa00add78eff6bf0a6982766a035bf0b465786632c5bb240788ca0fdf032d8815899353ea4bec5848fd30118711e5b356bde8a0da074cc25709623225e734ff5bd0cf65c40d9fd8fccf746d8f8f35145bcebcf378d2b086e57d78b11e84f47fa467c4d037f92bff6dd4e934e0189b58193f24c4222ffb72b5c06361cf68ca64345bc3e230cc0f40063ad5f45b1659c643662996328c2eeddcd760d6f7c9cbae081ccc065844f7ea78c858564a408979764de882793706acc67d88092790dff567ed914b03355330932616a0f26f994b963791f0b1dbd8df979db86d1ea490700a3120293c3c2b10bef10203010001a33c303a300e0603551d0f0101ff0404030204b030130603551d25040c300a06082b0601050507030230130603551d25040c300a06082b06010505070302300d06092a864886f70d01010b05000382020100a21a0d73ec9ef4eb39f810557ac70b0b775772b8bae5f42c98565bc50b5b2c57317aa9cb1da12f55d0aac7bb36a00cd4fd0d7384c4efa284b53520c5a3c4b8a65240b393eeab02c802ea146c0728c3481c9e8d3aaad9d4dd7607103dcfaa96da83460adbe18174ed5b71bde7b0a93d4fb52234a9ff54e3fd25c5b74790dfb090f2e59dc5907357f510cc3a0b70ccdb87aee214def794b316224f318b471ffa13b66e44b467670e881cb1628c99c048a503376d9b6d7b8eef2e7be47ff7d5c1d56221f4cf7fa2519b594cb5917815c64dc75d8d281bcc99b5a12899b08f2ca0f189857b64a1afc5963337f3dd6e79390e85221569f6dbbb13aadce06a3dfb5032f0cc454809627872cd7cd0cea5eba187723f07652c8abc3fc42bd62136fc66287f2cc19a7cb416923ad1862d7f820b55cacb65e43731cb6df780e2651e457a3438456aeeeb278ad9c0ad2e760f6c1cbe276eeb621c8a4e609b5f2d902beb3212e3e45df99497021ff536d0b56390c5d785a8bf7909f6b61bdc705d7d92ae22f58e7b075f164a0450d82d8286bf449072751636ab5185f59f518b845a75d112d6f7b65223479202cff67635e2ad88106bc8a0cc9352d87c5b182ac19a4680a958d814a093acf46730f87da0df6926291d02590f215041b44a0a1a32eeb3a52cddabc3d256689bace18a8d85e644cf9137cce3718f7caac1cb16ae06e874f4c701000000010300000200b8e3a4d9a4394bac714cb597f5ba1a81865185e35c782d0317e7abc0b52d49ff8e10f787bedf86f08148e3dbd2d2d478caa2a2893d31db7d5ee51339883fe84d3004440f16cb3797a7fab0f627d3ebd79217e995488e785cd6bb7b96b9d306f8109daa9cfc4162f9839f60fb965bcb3b56a5fa787549c153a4c80027398f73a617b90b7f24f437b140cd3ac832c0b75ec98b9423b275782988a9fd426937b8f82fbb0e88a622934643fb6335c1a080a4d13125544b04585d5f5295be7cd2c8be364246ea3d5df3e837b39a85074575a1fa2f4799050460110bdfb20795c8a9172a20f61b95e1c5c43eccd0c2c155b67385366142c63409cb3fb488e7aba6c8930f7f151abf1c24a54bd21c3f7a06856ea9db35beddecb30d2c61f533a3d0590bdbb438c6f2a2286dfc3c71b383354f0abad72771c2cc3687b50c2298783e53857cf26058ed78d0c1cf53786eb8d006a058ee3c85a7b2b836b5d03ef782709ce8f2725548e557b3de45a395a669a15f1d910e97015d22ac70020cab7e2531e8b1f739b023b49e742203e9e19a7fe0053826a9a2fe2e118d3b83498c2cb308573202ad41aa4a390aee4b6b5dd2164e5c5cd1b5f68b7d5632cf7dbb9a9139663c9aac53a74b2c6fc73cad80e228a186ba027f6f32f0182d62503e04fcced385f2e7d2e11c00940622ebd533b4d144689082f9777e5b16c36f9af9066e0ad6564d43"
	blockBytes, err := hex.DecodeString(blockHex)
	require.NoError(err)

	_, err = Parse(blockBytes)
	require.ErrorIs(err, errInvalidCertificate)
}

func TestParseLargePublicKey(t *testing.T) {
	require := require.New(t)

	blockHex := "0000000000000100000000000000000000000000000000000000000000000000000000000000000000000000007b0000000000000002000010a13082109d30820885a003020102020100300d06092a864886f70d01010b050030003020170d3939313233313030303030305a180f32313233303830313232303335335a300030820822300d06092a864886f70d01010105000382080f003082080a0282080100c0b765254922fbab6627f66e1488ed193180a9ba4f03720bca63d8726ace390214a07ab354c3138b58013d03c9401e6ca298d04922faf04d3b675affc99a50ab102ef2be67bf7dfa0dff2891597206b51293744a425c3a7392d2fc87eb6f8fc95421dab50144bd67165cf50f5e9d30f4d3239d780c4b0eb318859eb8a0b2f9446cf370f9a573190f9a3e147ac882e387f920e75081c92ee88c7ae52202184f8d46876305b416bdd70e9a28615a3eaea0955912495ec1e2deea938fc248b8dd05ceb63452e21e950a06cfc987e980ca74df4320d34682980388c513559b10b9da884fd7c3309904d48a2ec7a8d907aa78bc37c7e1c40cf735100a5540e5ff01c4f282cd9680685f722551e270d152bf590aea5f27592940b5bd2d927b57c817bccad4c54666e2a06ec7d18c56a252794f91ac1837ee7500a664b1e716da23bbff137f684bea18c0d3d0350b1d66f3e8d39e61bef5e515b338f9202c206453ba59fafe4250e7c4bdb0e190632e5c43513208c367b75bc0d87e833460f842c35119cb550d6105ca8fa421749f28cf20e14be0457350ced3681340107ddd540c73e7e35329a82b1bbea9f89da7518327207fdac5348d89ba75b0c2250c3d71c1c38f9656b9f4f557adbd3ad294fe64a177c86a328c88fdd1e9ebca7802326a4a470e0c96a6cc3f434d57ece83e06fe43590a1cebc91574f8013350b157cc52a3975b9aa6d5102d0c0ac7acfca0a14fda6575888998346397013e723e87e6b12b376b602bde38adb005227572be1e0f0d9b7ccbbda95c5447aa4d82ddb1e008025b6082c83a8a6ec597c7eddf2b77dd436ea19d76fdca9d8cbed9053c3107c6f190261acd28cdbe4f1844573bae95a70b9d6943401b0634802a2f8e932432b653608b2010d11c7252af66c85b86f866bf039e242fc0176f9d0dc8b911c0512255bbc7ffa91a1c1f770011ed9f83c3ea21aa8812c43d83c7e478e0af65c0a3ccea59b89aca74b42b38f6bac9d63356d84fd3c0e43c8d417df9bcf23cb8c23f0ba94a98d8703a87e2bb9f172052161ccacac35ae0666931c05b4e1dd75c69c24c9a9b70d1333c0cfda76591342a0016ea151c1fe92b4d2d8bc05fe7520d2a71d794c5de07358becb46dd52f9c87d7711ef892a5f436ab088664e0474da87500d0bbca383f0203a4d3ab9b0459cc8c52a9573ed9c30891a6abc6abf59bfed6f4f57f4aba63fa8be56280893f01d958925b405d9e18fc8dfbccfc53ebed9462194953ad923e656e9ccd2a5952eb887b65a8b2dcbb47c9e4f1d9910834bba3518eeb75169c5ea1a11954f5ba5087292cc529218cd7ac560a31fa1bc0a912bbce45a7ec85f5b74709369b7269507b069b0df6439eddae6c241d72e685866b35a80bd67b68934b8efeb5079afce470f48221b47acccf3aa7d051e5c02e8125daad22aed88e99c33a217c90912ba81ce6ad884b34aa82fc300dea6d418c1a07417909dbfe5e9a7649b5b95705a329e2d29be70ce73306cba0a9e3750bb214c6536b605a7ee2e820018f338efe13c47c38d8531dfce16efa9b17486d84c8ac6217df461b4d00058279e5e877ddf042060161b1496754825e5f10e9a962334a767ddf237832834239d28f6ee1413b65519306fc3dd2b9f5b36af3d7beefece921415ff053c7738a6360be21d5234e823fba39fa1de8dd8b30af66084c5be5cb1f7c11d7d84a9b41e9c27d2213e00c6bfe7e3498b08b2ec6fe60d8c7c4f662f48cc0c7bc1edf35b63be27f6d35a41288cdeb67fa996e3b3b0056f8ac162b0d8ba7aabaff05b5159a1afa2553bf14f316806e340e570dbe6b28328d2d8ebd69a807e4a3a4410bf1a84c8ecf8ca74f822691815cae3fbea424a88bb7dd038de1db7dd6959acb6d89a7e98238be2ffa661f8a82193c0bdd98c16617164bee2032397369a2bbadc2c4c8b202d41d5da374eef83c4cab03e2d0fb974c1c60e2c0e53cffdc584586e0d74e8328d6cc19d44b75e27ef0413a2a627aae63284a0817a42d6292a3cfa34578e3cd2206ef31f30177fba651708a5c8a50ac22ae2106c8324aa92de33b99a29a3216f2017cb892cbde504eb6ff8b8ab5792e7e940be758a2e584fde782779daffcbfab6104adb82f9e32cdc507bedd745e491df8c29cb2ef8bf1de78e0921aa19830cd71a8499615a482e73cbfad2b62f25c010f5bba77270b0cb8061bcaeed6f50007fdd784940744b561a7833442c32f5bb5b644824178f008881b71a567063b0d8563cae02f52a2a572f050b369cac71525ed548336322e7723e56e8a0c85b7b7247c982e2464a1344dff3bdb62f58c1ab79b1d7d30ba2ff2d54d33f016594a8bb16f98ec3c01309d7d9588af8678b37220b8c5122a66be6cad0e50e6e51cb20e77463818da59f8fdac6d326a3a519991dd51026ba4bfb6d40ef3f8e7278aa40b4a13fe90bb31104c2b9b6c6b474cf2a3c1abe3126d20b938df64a87f796a9ae1b5279e595a3b257f556fc577885a7fd0fa940b1aaf8123077122a2d5b81a3ee80cbca5d4d26abde55a51d0c719b76ccbd252eb80464bae93b0563c405a9bf23a86b489ca6332cd0f328599a1ebee171b76b1680d34a022ba9ba6832733ed754a3fe21a094de2a7917739467ec6db16c1c39535189e2b4174bf03375abb5071b04cee706a004f07a962a95a8106e281813d869d7bcf315dcd44a70948332013aff82b1375b5c56e52b080dab625be87040937ddd74d61cf97b199736d44eca1910ef47d77e4e30e2ae6ee40b8b4dad07fc49e0fc0f7cff80ec809da173511418b0e8d72037750dee55adcb94254d43fe042c37c7a9a480b8026f9b241f47fc9c4d64bb25c146773fdf1ee60131efd82ac94966c81a80e1dfbacb752ca508f110203010001a320301e300e0603551d0f0101ff0404030204b0300c0603551d130101ff04023000300d06092a864886f70d01010b050003820801003c0c9f58c31bd4e4106ed260ade850926c60ac3a932d4fdd81df0b6669fececeaa0f0e17d2b3d6f32acdc43e343d983cd65280332b7345706ac74c96ba3b94eadfdf6106de13c2053cf98d817075e4fd9a5dc3114a9de76f68222efdb9f738420a113801979c3a0c4fd560c06ec39b85155909dad50e9c84ebff12f0e159fdcf2d6fe6430ed07888f798ffdfa955854bc13361a984735122d833347696ad4459f918fffafc10b566f830b8e3ca2c965ef662c99e4eb38eb82573d7d71618f79e055941cc3777bf1968d09e2635a029714f8572930b900fcc5538ed14d56b19965b9df1d5ce487d30decd8c4efc59d9161072cfb6bb88226ac8029a930ae0f46ae8e0aafed014ba60a43c0fd19f62ceaf170db8e64c885ffd7ab9141691fe477e15b229597ff072ce22e3ffc92d2af59799b151ddc1fb5b113d5917f048fab339480c7ac4daf0615edc39737efa07ea82e50fb2dea2edd33d2e5fa413320de814db78a2dd180680ecfd1f62ef5279e3d4bca4c298420c5067f51d6fcdbc3e3ad21a4073b2b92f3c747979e439f3e426b41dd9e692a739928b813a20f3366071a160cb401802fe25247339a1669cdbf7de5f4b0047fc023fbf1f286275db74cf1e3c53ba753bb7f01694df73b60f9c43bd39f3e2db7d56d762a0161ec2c6e63a383747be7e59dec28a418ece4d864579757bcb4e864498743c50b53e28c27cc0269a3bc8b2306d54a05ae81885c7e7e2e55203b7665a84800e91b37db0191a237cf30cfa218a4e2339fcccc8e9a64a9dd2e3b81364fe49c022cc0343e3dadb92edb6a49f3d382abf8eec787065068cf98ef252d50931894af3601ae395e965d056fd3d763f03f17b84af26749c731b6c8474517060d7b6cd7b6f8efadb856a302d24e1e0c09cb1bae5278d1f2bb11ab9653d3995f0e9412863787106fc237e1778d8f1ae0c13ba0025988a3628f172e907ed85fcd7f2cce8996c04ec7cdac22f8efafa44b8d18a34aa61e209bf937402a5a4fd4af46d312c62827635d05ab6e8b45862aff3ff72da21dc30149ca1881e8e7faf891111f79469df89a7237c3294b5b8111beae5fdbb07afdfc7b47ffc6a860b202eeb7319d4d5669a29b4071e4a875180ce8ee00175efc07e0a0eaa3465e4a75befec33fa54a11568ed05942ca779d9e103bc9676c2931aaf38cb8b9c9bd0ffe05cfd471b5544824f59e7f701c303bcecf0ba64bbe696878e4681376be0fe404ff056b340384dd886a6e7fd7566f1c9ec1e491c0da28fdf7091407a4099a3a7399682f6b4b389fab73adeee912935d138ce259ea40720588550e3fc08497d1ee62cb85e3b2ddc85728af82c86b1150067b8a2d2c857b1a23c951b9f326bf8e22fec2b3f8ae1f8fc06e25ca34ed2ec6d210db17bf2230df79ea7deeda3039f8831a18b148d83c0076af887275d607ebb884eeda770dd5721884a625d5a7035bda08c5756b31c6cfc3a9b547603dfd8668a488046336096603476b7a5e19a053b486b49042baf7dab3d45b3fc15ef443821754c0c42bb4584a1e47aafeac1730fe4e075e7bcc6c82c7a637c642da6a3058fb31c6f0064358bfac21798806dd96ab358588b127b0a3557af4e9b6c7342d31a03319c64ddba260ad2b210112686ea104e25262d8cdd5dbfc991f5545af11299ea2371f53b13fed0231954cfea289b03b0de15b5782d03e785f7388b3e53d5fe810fd3a520519813167a67c0f09ef95357740865eb4c73561098aa9802dbbbb038f6e462eda0b6667f2b26fa053a4471de3bcb0901028e7a44d7e791f1425e2324d6b6d15a40e38d92ebb36cd63f60cc16a4046b09919b744283c2c7b0d1355dedd131cdd8f3a49bc6d3770372ef70affefd98e74a11605946e868bef4057853d50bcf76b7b8742c9a038c700cf9cd8919058e2cecd66cdd55f3d71290471f1c92a44fa9db26ae23b06627bd4dc4dfcebfa89d5e0fd0186772723005c90046d0acbc54a7c9567fa6b9cf3a06145c7e73d53da7aa9904f625ffc648d4a6a27d42b09efac2638ca4ef434278a6e4fc0a8b5b39e229e7fa80714b6272c7910e32bf24a73149d1a4a64851e7f52b81ee7a5191ed8ff2070b7080ed4211ce6da06a5b2af5d4a5ebafee67bb8b6e0032332f1452593440a8e1e7ba0bf032740ca856edd68653c62c9958ea6db839aa4efb15aad96711dd8928a657e15a3a6c713f78919d33d0a9f53c6f00886e085395b307bcd90ee3c799d452c6eb14ce535377327b9da505c7231e7cc067a4036589ddae912a857746cd7f03270fb526b111efe4fdb95dda2cbd9b8ede7f05bff68a1bd13e831ac2222a21df1acfd3d50d9701fc96c427975ed012936cb5b15e2cbee1fbc0fec92af7d91d38cc75a8815a206ca5e1b7067df070ff64635aff10966ecedcab095e24926e41807f91395ee3b54053d092fff71d6224f9e0d1d08fa3c4b6b16d050964efae26cae9ccd6355ee98fa9930f2da87d52154309590340041d449b4e2e2127b09d067f2c012f67d295f12bd1a4990eabbae2b989f88c56ebeb9290ff59e35626cd0553901cd2cce01a4d1ca4fd1231ab583d429fb96706932af90ebee5a10b8cf356c4c8590f7d793af28d74471acdfbb941b6bfab3faddc632b261f0abfe67ffcb66b63c289c7876fe58061d92d74aef3bb23d8520298244aad891b74e432c3d7c1cfcd5b4cdd08f0e0d149da47f653988105da4f58cb39d3b18844773052ea4a1b54ec999b4fdc930f42a54fd6c412f69b789579d7861aa847b6c59868f7d0e8b54800cbf0cd036ef747b805f5e13bc6a65e17dd57eeb5535cab162a4fff7261aba35f5a57b0a8fe7a0cd2e025637727a3c261a3025a44f8dff30135d02635b25d4643a42caa0915c7000000010300000800b66be4e0077673a6a8abcb3e323eb06ee44d1cfb9b45679e82bd6211a371ee919a4fa3b28cf9d58f3673703abf3962faba390b250a0879f947ff777f2ddba1ede98e087ad93abdde58edd64def32eda4898a24aca83f7441eedcc8bd16d36bd2c54b01ac996b54bd2a14e310fea353ff609c459a0a0c564dba51b8d08cd2909312c14c82251a00c5bbfe0884a6bf71997b7152559dfa0b4935b21628790f375ba13213a7b25ccdd1876d9f545a839ed036be01af778efdb451201856037bbc7b736ddaefc3d13082a9422d58b4f11ad6985c4f365763e2ef35453f417b53cde34d96ced798145b15e3f94d7fab214bac1e331b5ed67e25520c6942cd0b02e09b5395a6ce84c34ab1c36912379bdd08f5a58c6abd98479cb4d5e87431fa5a71f4820121d5cef9c3503c9d79da56be3a67ee81c9e8c0b3dbeb76430c46b09c6401fbbf296868d05c5ee5b1e10688e8aabb7bc6d8c915c4b6346380b6b4a1a95a8c16eef2dfc7853bbbb81360e0840ef69df09ff2cc34f0c27fef21126e29402f24e3ad39ea1aaff52efc837bf8851ab977bcb0bc0875e6ff2cf3d65fa1abd30bef621dade8d4b333f478c8e00321a8282c071b24703513a4f0acfa8b0c98022f09460c2b13f89ecd7092d47c399cc2722e814104db23cacddce4fe0685d3e94e6ac668410f2ac94273a393f9b25a9e3dae3af3d2a07dfd8d90915fc01b16cd649cd45fa90bfed1532afe7f4834862a6fd3dcc086ddace177924ef57a61ea0d6e9fc4c2dba0ab9c847b3c376ec4d24491e3424a1f7f1171eadcad1d8352627ebc59f0af4e9b88f33375f742dfa170316ee2d3c7fa22abbdc50010736fd34784b7eece6d6c18bdeaaeae5df44e3847ea143ddb4da3d377844d60f7159a43f175b6d5765951d04045cf619657e683aa172a59125df40f85b81f232077cff06d5f8e34fd94818773b71fdef0a8f9fe364dca33645a7432b68ba44d2c872f7421f7288552cac1d1233715f83e60029b3549cd588399d9b107ebb3d7d1cfb762886483b4cbdd8222c481ba65057ad68a39d9b0f6a5f358c2482b5d2786b9f9ca615c096d82e345eceb3333746ab60a07a72331cc4f6e58c9126a7229d162084a58e7273bdff1b6accdf456a6d82d9cb7b2b8683191f81f41dda147db727f1013420be7f737fb64c4e4c53d3fc5afd794f7dbd14b5ae935ce1eda7753c78a08627d75e5b54e089e5696dafe5efe037a6452616c57954a7056efeeb6712858d70b8c52422e955281443308987e0e0a53d1fa2cd1ae4a65ce384495f01c100e253b653f71503cce5da0958b397323ca1462f2072e0c19f8de1b9ff9e7e43bc909ebade3cc6b4b1a9249ad0dc8ba2a1c98c3775ce2cae765adcdc4ab4b7234f7d7fd9fbd8243a2bf908fd605c59eaf394f1f9c945220a021f04df6f69c724f8f01cbd3349185ca254491c57207e30bd6c8c6b47a5aaa0dd52e7e713bcdcfd1da186c2d7cdd7c7cd3fb60b62c97a18c34127725f7c8cbea79da03067facf364e5fae7ea0b16474d21664c82ffd7b63a45896472c0c540335e0a1310e252fe51adf0489c9cd5162b723f680f01b8ce4904989ca85f0c171512cf6b70800f6612cae4755eaa47319622f08a113a5926adbb55e2a1c076bc61510f42a26ba33b12f02952e904c36eb311a908979a66fa188835543a6c61b70ba2a80c1d3a665a2b8758f10c7584110369d5e6414b9d092afbf59c875d0a3b448353e75b158b320e853d8a57b463e42f3c766099f05d4dab682c26068d4303bca9f19fb88ab3c9ea090aeac15cf7a0eea22e999abd5bdb5ac2a9618c1fff4dec60ea67c787d8feb37d3f08d8e7a2cd7b943977460b60049863a106f036efe3fa8f1ca4f11895fcfc81599d6add7bc9dd98c8fc537e8c4a006a7ad1a9ae8d7de6d17da7533326dd1ad183b123df8941b9750c5b8371f0fd6cdbe23b27cd4d75ff314b9c77140cab5865618e32d3ca82db8ab4ae9b8173d83e91a8f455ba50fdedd6b975836ca2d1a2cedd98266e4c739b4d76a9fb7d97ab23cf29cc90a9ce930f66dd52d8cbc9759f1f35bdd9ba0ac92929613149d15352cf0bf9d40b159df6ea9f2739ea82347fd1f71b94d6d323a978a03a0546980ab032f23aadacb0ce27e64f994103264e654f0521490fa12f55016189fe83be393dae97064a393cd33da8e9f5ebbe85d10532f476f88a3f8923900babae985a1d4db922048ce6d81ddedb06c93982afdb5141f72519ff053da266210e94d769506c5da48c230f0e89faddb9ee002f46f9668cb4a649a4be0308f0720d77d9bd5c0c83a2b302a3c39ab6c6e73c9257671a21b9a6db0743986d733edcad29a2c0ca0dec8aa310684ba1d92b0ea3b62c14aacd18947c60de08f544fc1ef93db6fe94dd878f36093d14eea4230e36e6e67492d1f046ada7efe526a9ebdbb57f06c0e4099c28030155ab3a4f788ed4ecf733cf6277adbfabeb5dd283c5c29392be4e7131d6f443f212e6d9087434b0ea68a735eb52676ce1de9a85e90fc8ef3992eb193be1b8a1783865b3755a3257e28cb255907889131a7304c768bcbdcff3e6771de1c47350664310d5ead91bf680c00df79c9b72a577429589e265d0f1c715dac96feafe20d8150e5cbfc441c5202d5af961176314d29a9a9415bcb25fdbbe1a7ddd6d0e737cdeaef6ad2b0e082b36cc4a3882110479e36591c51b693d127f6e2579783fb4d65bfe42a4b944a489c08485e8589c7d8bd6a9a81b580254760cb3a64130d46ed28234fdf7ddea2e7a0e397776292f933424ef3e789306d85856e3eff2c0b3806360bec39072d20cbc00b9c57f396e7f4f99b6a1753c9128e0568a28b0a3c62443813755f1e189be30b45872a46dd963c3921dae2d"
	blockBytes, err := hex.DecodeString(blockHex)
	require.NoError(err)

	_, err = Parse(blockBytes)
	require.ErrorIs(err, errInvalidPublicKey)
}

func TestParseHeader(t *testing.T) {
	require := require.New(t)

	chainID := ids.ID{1}
	parentID := ids.ID{2}
	bodyID := ids.ID{3}

	builtHeader, err := BuildHeader(
		chainID,
		parentID,
		bodyID,
	)
	require.NoError(err)

	builtHeaderBytes := builtHeader.Bytes()

	parsedHeader, err := ParseHeader(builtHeaderBytes)
	require.NoError(err)

	equalHeader(require, builtHeader, parsedHeader)
}

func TestParseOption(t *testing.T) {
	require := require.New(t)

	parentID := ids.ID{1}
	innerBlockBytes := []byte{3}

	builtOption, err := BuildOption(parentID, innerBlockBytes)
	require.NoError(err)

	builtOptionBytes := builtOption.Bytes()

	parsedOption, err := Parse(builtOptionBytes)
	require.NoError(err)

	equalOption(require, builtOption, parsedOption)
}

func TestParseUnsigned(t *testing.T) {
	require := require.New(t)

	parentID := ids.ID{1}
	timestamp := time.Unix(123, 0)
	pChainHeight := uint64(2)
	innerBlockBytes := []byte{3}

	builtBlock, err := BuildUnsigned(parentID, timestamp, pChainHeight, innerBlockBytes)
	require.NoError(err)

	builtBlockBytes := builtBlock.Bytes()

	parsedBlockIntf, err := Parse(builtBlockBytes)
	require.NoError(err)

	parsedBlock, ok := parsedBlockIntf.(SignedBlock)
	require.True(ok)

	equal(require, ids.Empty, builtBlock, parsedBlock)
}

func TestParseGibberish(t *testing.T) {
	require := require.New(t)

	bytes := []byte{0, 1, 2, 3, 4, 5}

	_, err := Parse(bytes)
	require.ErrorIs(err, codec.ErrUnknownVersion)
}
