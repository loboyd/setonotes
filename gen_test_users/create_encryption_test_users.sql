DO $$
DECLARE ne_test_user_id INTEGER;
DECLARE ne_test_page_id INTEGER;
BEGIN
    -- delete existing user under same username
    DELETE FROM users WHERE username='ne_test';

    -- create a test user under the new encryption scheme
    INSERT INTO users (
        username,
        email,
        password_hash,
        main_key_encrypted,
        salt,
        version)
    VALUES (
        'ne_test',
        'ne_test@setonotes.com',
        '$2a$04$JfbLIUmgJ0886T3HJlrPQO6uoBzatYeqJJuoc7YXd1eo2Wm59f7RG',
        '\x42c1b27beeadb4ca6c195141f79ce69ac3cecba0bd74e491ea4324f4f0e7fc91e1ebf8c81fac45ae0da5886a',
        '\xfb5a7663386636729fa8ec9078feb7e6',
        2)
    RETURNING id INTO ne_test_user_id;

    RAISE NOTICE 'ne_test user_id: %', ne_test_user_id;

    -- create a page for the new-encryption user
    INSERT INTO pages (
        title,
        body,
        author_id,
        version)
    VALUES (
        '\x917e9fd42f3e57736e7f9a6a905daf1c0dec4cf4682e863fca1a05d506fcd22c1cb330bd7a040974ec215646051ffd849e9dcd116acf',
        '\xcf0aa0f55c5ba2113023b7daa973221a801baa635e97d4d46e3c16dbfd3b14482be26e517c4ea48ab58f270e0b232e5e88d7584c17c46bb61b01a67947e4b2798a866be0d78349816c9a046582473515c71024bb26b033e4209a79efe98b32ed51c717f799e8043c8f4c3d37310e91e028158770a8aba3ffb037045087ba536a944ce0b4630a7e6a639eb5df1807103d08bc112f3888933d95d1ade1c235da71929645a0975bb2709bec9b9f701356e605b32acba839989ddfec7ba99ec4691bba3dd7df8f2547873106894cd09d39cdb691807b5c5d414f457957023d7a197f46f8b85c9ca660cd8e4d87fb9f9df5a736676fc4c58329a31a6f53b1ca30a34194c8e679963ea7f8ccd95817c31407fa0c067db929d172323412ed5e72de8075e1e79f10834ecd67bdc671fc5d2d642e2b82e2d2df21c6e9743024a7990105bfaeb83bc2fbc62d5e7312060bd97dd172b2fd117e1a8456787b79ee22dc6ea73c597f644bbdef88d648a1057156b26917ff4cc02d1897677b102e29d3663aafcce432de5fc9706bb6e82a695ec57174440143e79718f694a28fb804b626fd2a7e5e3e500dbea809e661946efe888b4d63afed56aaa98fce463cd444cfda4b8ae49fd116b68f3a0622bb380805f87515c298c24bff17d420d9753234127eab9ea08b527abc6d1aef7f35acf74852e2c445eb449605641b0bfa9141d908eeb99e1d4d507297a26fa3704939a72cea767048f12870ab1c2af6e1eaa4980693cccdac2393afab7f9f5e5b3606f7d19d5509e10036f039b7e4c39df07f3fbd5dabf7479b95aa7a45665f6ac2c62ad13c8de1e236fe6a6c226189f1138700c53de4b413c9857f17370bc2c59d008a12ccd376d20e82e5bfa8045c893d8b175aefbca554ae722276e2917dd722f2f3b05bf008b7cf8ebdab97aea392eb2dcbcf56e93fcb26e4725ecd62e345ac81b50a1ae106a1290f72fd1b45488ecda7d0a1a5b86df24434ea13e4e254f9401c966106d6d0e238ef0145b82841c7feabd0689e33faabcb4724fd7a4c28bd44e95e5d03967ee040413cd68b45c6ed33a45e40925e3e74b0518a8e2c530821b5b73a8bff6b01739a3cc0537b8de6a5de5ea1645b57dcbb61ba3c88781c3bc1a63da4b66c541b08',
        ne_test_user_id,
        2)
    RETURNING id INTO ne_test_page_id;

    RAISE NOTICE 'ne_test page_id: %', ne_test_page_id;

    -- create page permission
    INSERT INTO page_permissions (
        user_id,
        page_id,
        is_owner,
        can_edit,
        user_encrypted_page_key)
    VALUES (
        ne_test_user_id,
        ne_test_page_id,
        true,
        true,
        '\x99c6975fc066a83e666623d3f4ed3ab21b5a8b206056b2b39c3a3e3d67a3a38faaa5761297681001d099cd16');
END$$;

DO $$
DECLARE oe_test_user_id INTEGER;
DECLARE oe_test_page_id INTEGER;
BEGIN
    -- delete existing user under same username
    DELETE FROM users WHERE username='oe_test';

    -- create a test user under the old encryption scheme
    INSERT INTO users (
        username,
        email,
        password_hash,
        version)
    VALUES (
        'oe_test',
        'oe_test@setonotes.com',
        '$2a$04$sabkxdQk/xPCuJ3SICOmt.MvNczcxAp6u6v9b2oBEtrYf66YWYeOe',
        1
        )
    RETURNING id INTO oe_test_user_id;

    RAISE NOTICE 'oe_test user_id: %', oe_test_user_id;

    -- create a page for the old-encryption user
    INSERT INTO pages (
        title,
        body,
        author_id,
        version)
    VALUES (
        '\xdc31701917ad86205483181370bf814b449f239b754e365f948c80c4c83b8b7f2fcb8c49ee6d2f3266567a5f6f699645ec5e0e3cc10f',
        '\xde02dfc82c45245d38c6949c099761b59cf1689a813d0f08c3216fb8a23305937b102f26269efaf209a6f0b82e4715342d6ef39e7e46ecb366f2001b75352f83cf2074bbd3c89fe55dc5c36137863f3190c871c766a6a156cc53157ea4744544a100ae9ce5bc965536ba3bef5d2217198537392e324d82fdad1dee4082ea33672b74672517efac3ccc210266b70e3625b86b1b99edca371e993103b23c61a3de7f963cb1f0f903a4c81a9fe2572a6942a4022b57de238273105a1c97c320186a8b3c7a66c61426b4f6f6353c6a5ddf9c3f82c9f855dbdddd6448ba5890f7882f6a741c5f26cdc14f079dd200f1cf34f3c8add816e06da0845b0fad7f92ed78f27d7d5d84add01829562ecceb853cf1e32ec750b835651b9b0babdf9b3f4e38288f0a0f52a52c7dd0238c7f148cdf7dcbb9cf63377ec1a7661808c05c90f885195bd1a8484a95e64ac0f320b7839bccd8f6f9abd478aa4b8b3a74773cd8c8982bfcb1d592ff5340fa401a5475572a9a26c49a9fcbf992c6fad5b427566304b430a01de6ad68eb0f2a22f7403afe259d2db41ea067a58d9c053165f499d1d5a7a3b3ca906e4739e1401cb422d660002547b7f736e0e5356d44a0e9865a9ea17f0784964bc557e58b83639223bc7a43050e8538985c0515d89efea6030862a947ed9a6a8ef3bc7653dc5457f4c55479999b85ce3b8c265ffbc61baff349c0e46db23aadc62ee82add36ed4dd9c8920ae6979bb61014970df9755f07df44f82d3e0cbede3d43bcef575e326a14579ab1bffd5fc732a7167dd1baab7124679a79667b5fb8a35b59447ff4bb590290d7cb1eeb2b2c7ca509b17364a965fa1bd1f035193a66e6dab67c33693ca5b56daad4619bd0dcd05670aa5e47c2a6dbd511d701f77f3165edbf1edd434c3835a1ae2bc49e359b6992b5b3e94e7bca8b1ca090bc31ef58a4eebc67a0155add18997ebf18c27abc596078773e28ac6e81da92891354ef3a6e9c5988589017f165802e98d93ec7f27ac40734a719efa1eddbea74764918707562f06f760d00b29a1c2492de74bc72cddfc5cbaecb22e86685c5dd778799b8bc864960bcfc323dab85bb7cbd4e1f3c014ddf80e115f729e9f57a2560f7bf1091b0cb8d883a0b1c75ebedc2ddd4',
        oe_test_user_id,
        1)
    RETURNING id INTO oe_test_page_id;

    RAISE NOTICE 'oe_test page_id: %', oe_test_page_id;

    -- create second page for the old-encryption user
    INSERT INTO pages (
        title,
        body,
        author_id,
        version)
    VALUES (
        '\xfbfc66192ab0ccf76ea2a4704ce0f53d9a9316ed12360b4af7f05024fdb3071c82032582f5dc48',
        '\x9e8d2c77c75f1ec542e9b0b8a8a13f287d50c57ba88c3d0dcc24fbd49cf8ac52ff43d1736fa9706d2e2944a38cffb4de83b3f46416',
        oe_test_user_id,
        1)
    RETURNING id INTO oe_test_page_id;

    RAISE NOTICE 'oe_test page_id: %', oe_test_page_id;

    -- create third page for the old-encryption user
    INSERT INTO pages (
        title,
        body,
        author_id,
        version)
    VALUES (
        '\xc5d70c2c65759d2b4075f3d97d37f10a8cb5be88846c7da015fb1cb2482f30babebf0bdaefe65bc4c3b7fa0b',
        '\xf3c4a634bac631dfbf7dfcbc6cc7a2c215a4820924a9f740c9b0cd34ae7a08da39cf94205d667ba0c67ad7b30122a1eaed',
        oe_test_user_id,
        1)
    RETURNING id INTO oe_test_page_id;

    RAISE NOTICE 'oe_test page_id: %', oe_test_page_id;

    -- create fourth page for the old-encryption user
    INSERT INTO pages (
        title,
        body,
        author_id,
        version)
    VALUES (
        '\x6d704a0d8049009b04d8980b1a3824589a2e2183119f84eb78046559fe38672bf928b3a2288a35fb825ab0',
        '\x7ee82ceedaa4bb67f70b56161b01b8df8f016ce6e9361f8eb8f442d260114de2070cd884ae4ab6630a84d7cbcdc83d445bafaed9cb6aa506874de15a3c36fd9c34b16ec079ebcf',
        oe_test_user_id,
        1)
    RETURNING id INTO oe_test_page_id;

    RAISE NOTICE 'oe_test page_id: %', oe_test_page_id;

END$$;
