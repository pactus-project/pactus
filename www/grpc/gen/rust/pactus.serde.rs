// @generated
impl serde::Serialize for AccountInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.hash.is_empty() {
            len += 1;
        }
        if !self.data.is_empty() {
            len += 1;
        }
        if self.number != 0 {
            len += 1;
        }
        if self.balance != 0 {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.AccountInfo", len)?;
        if !self.hash.is_empty() {
            struct_ser.serialize_field("hash", &self.hash)?;
        }
        if !self.data.is_empty() {
            struct_ser.serialize_field("data", &self.data)?;
        }
        if self.number != 0 {
            struct_ser.serialize_field("number", &self.number)?;
        }
        if self.balance != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("balance", ToString::to_string(&self.balance).as_str())?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for AccountInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "hash",
            "data",
            "number",
            "balance",
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Hash,
            Data,
            Number,
            Balance,
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "hash" => Ok(GeneratedField::Hash),
                            "data" => Ok(GeneratedField::Data),
                            "number" => Ok(GeneratedField::Number),
                            "balance" => Ok(GeneratedField::Balance),
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AccountInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.AccountInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<AccountInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut hash__ = None;
                let mut data__ = None;
                let mut number__ = None;
                let mut balance__ = None;
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Hash => {
                            if hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hash"));
                            }
                            hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Data => {
                            if data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("data"));
                            }
                            data__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Balance => {
                            if balance__.is_some() {
                                return Err(serde::de::Error::duplicate_field("balance"));
                            }
                            balance__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(AccountInfo {
                    hash: hash__.unwrap_or_default(),
                    data: data__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    balance: balance__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.AccountInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for AddressInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        if !self.label.is_empty() {
            len += 1;
        }
        if !self.path.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.AddressInfo", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        if !self.label.is_empty() {
            struct_ser.serialize_field("label", &self.label)?;
        }
        if !self.path.is_empty() {
            struct_ser.serialize_field("path", &self.path)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for AddressInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
            "public_key",
            "publicKey",
            "label",
            "path",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
            PublicKey,
            Label,
            Path,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            "label" => Ok(GeneratedField::Label),
                            "path" => Ok(GeneratedField::Path),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AddressInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.AddressInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<AddressInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                let mut public_key__ = None;
                let mut label__ = None;
                let mut path__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Label => {
                            if label__.is_some() {
                                return Err(serde::de::Error::duplicate_field("label"));
                            }
                            label__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Path => {
                            if path__.is_some() {
                                return Err(serde::de::Error::duplicate_field("path"));
                            }
                            path__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(AddressInfo {
                    address: address__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                    label: label__.unwrap_or_default(),
                    path: path__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.AddressInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for AddressType {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Treasury => "ADDRESS_TYPE_TREASURY",
            Self::Validator => "ADDRESS_TYPE_VALIDATOR",
            Self::BlsAccount => "ADDRESS_TYPE_BLS_ACCOUNT",
            Self::Ed25519Account => "ADDRESS_TYPE_ED25519_ACCOUNT",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for AddressType {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ADDRESS_TYPE_TREASURY",
            "ADDRESS_TYPE_VALIDATOR",
            "ADDRESS_TYPE_BLS_ACCOUNT",
            "ADDRESS_TYPE_ED25519_ACCOUNT",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AddressType;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "ADDRESS_TYPE_TREASURY" => Ok(AddressType::Treasury),
                    "ADDRESS_TYPE_VALIDATOR" => Ok(AddressType::Validator),
                    "ADDRESS_TYPE_BLS_ACCOUNT" => Ok(AddressType::BlsAccount),
                    "ADDRESS_TYPE_ED25519_ACCOUNT" => Ok(AddressType::Ed25519Account),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for BlockHeaderInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.version != 0 {
            len += 1;
        }
        if !self.prev_block_hash.is_empty() {
            len += 1;
        }
        if !self.state_root.is_empty() {
            len += 1;
        }
        if !self.sortition_seed.is_empty() {
            len += 1;
        }
        if !self.proposer_address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.BlockHeaderInfo", len)?;
        if self.version != 0 {
            struct_ser.serialize_field("version", &self.version)?;
        }
        if !self.prev_block_hash.is_empty() {
            struct_ser.serialize_field("prevBlockHash", &self.prev_block_hash)?;
        }
        if !self.state_root.is_empty() {
            struct_ser.serialize_field("stateRoot", &self.state_root)?;
        }
        if !self.sortition_seed.is_empty() {
            struct_ser.serialize_field("sortitionSeed", &self.sortition_seed)?;
        }
        if !self.proposer_address.is_empty() {
            struct_ser.serialize_field("proposerAddress", &self.proposer_address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for BlockHeaderInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "version",
            "prev_block_hash",
            "prevBlockHash",
            "state_root",
            "stateRoot",
            "sortition_seed",
            "sortitionSeed",
            "proposer_address",
            "proposerAddress",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Version,
            PrevBlockHash,
            StateRoot,
            SortitionSeed,
            ProposerAddress,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "version" => Ok(GeneratedField::Version),
                            "prevBlockHash" | "prev_block_hash" => Ok(GeneratedField::PrevBlockHash),
                            "stateRoot" | "state_root" => Ok(GeneratedField::StateRoot),
                            "sortitionSeed" | "sortition_seed" => Ok(GeneratedField::SortitionSeed),
                            "proposerAddress" | "proposer_address" => Ok(GeneratedField::ProposerAddress),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BlockHeaderInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.BlockHeaderInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<BlockHeaderInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut version__ = None;
                let mut prev_block_hash__ = None;
                let mut state_root__ = None;
                let mut sortition_seed__ = None;
                let mut proposer_address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PrevBlockHash => {
                            if prev_block_hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("prevBlockHash"));
                            }
                            prev_block_hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::StateRoot => {
                            if state_root__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stateRoot"));
                            }
                            state_root__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SortitionSeed => {
                            if sortition_seed__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sortitionSeed"));
                            }
                            sortition_seed__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ProposerAddress => {
                            if proposer_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("proposerAddress"));
                            }
                            proposer_address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(BlockHeaderInfo {
                    version: version__.unwrap_or_default(),
                    prev_block_hash: prev_block_hash__.unwrap_or_default(),
                    state_root: state_root__.unwrap_or_default(),
                    sortition_seed: sortition_seed__.unwrap_or_default(),
                    proposer_address: proposer_address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.BlockHeaderInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for BlockVerbosity {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Data => "BLOCK_VERBOSITY_DATA",
            Self::Info => "BLOCK_VERBOSITY_INFO",
            Self::Transactions => "BLOCK_VERBOSITY_TRANSACTIONS",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for BlockVerbosity {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "BLOCK_VERBOSITY_DATA",
            "BLOCK_VERBOSITY_INFO",
            "BLOCK_VERBOSITY_TRANSACTIONS",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BlockVerbosity;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "BLOCK_VERBOSITY_DATA" => Ok(BlockVerbosity::Data),
                    "BLOCK_VERBOSITY_INFO" => Ok(BlockVerbosity::Info),
                    "BLOCK_VERBOSITY_TRANSACTIONS" => Ok(BlockVerbosity::Transactions),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for BroadcastTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.signed_raw_transaction.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.BroadcastTransactionRequest", len)?;
        if !self.signed_raw_transaction.is_empty() {
            struct_ser.serialize_field("signedRawTransaction", &self.signed_raw_transaction)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for BroadcastTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "signed_raw_transaction",
            "signedRawTransaction",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SignedRawTransaction,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "signedRawTransaction" | "signed_raw_transaction" => Ok(GeneratedField::SignedRawTransaction),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BroadcastTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.BroadcastTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<BroadcastTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut signed_raw_transaction__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SignedRawTransaction => {
                            if signed_raw_transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signedRawTransaction"));
                            }
                            signed_raw_transaction__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(BroadcastTransactionRequest {
                    signed_raw_transaction: signed_raw_transaction__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.BroadcastTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for BroadcastTransactionResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.BroadcastTransactionResponse", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for BroadcastTransactionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "id" => Ok(GeneratedField::Id),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BroadcastTransactionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.BroadcastTransactionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<BroadcastTransactionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(BroadcastTransactionResponse {
                    id: id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.BroadcastTransactionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CalculateFeeRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.amount != 0 {
            len += 1;
        }
        if self.payload_type != 0 {
            len += 1;
        }
        if self.fixed_amount {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.CalculateFeeRequest", len)?;
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        if self.payload_type != 0 {
            let v = PayloadType::try_from(self.payload_type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.payload_type)))?;
            struct_ser.serialize_field("payloadType", &v)?;
        }
        if self.fixed_amount {
            struct_ser.serialize_field("fixedAmount", &self.fixed_amount)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CalculateFeeRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "amount",
            "payload_type",
            "payloadType",
            "fixed_amount",
            "fixedAmount",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Amount,
            PayloadType,
            FixedAmount,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "amount" => Ok(GeneratedField::Amount),
                            "payloadType" | "payload_type" => Ok(GeneratedField::PayloadType),
                            "fixedAmount" | "fixed_amount" => Ok(GeneratedField::FixedAmount),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CalculateFeeRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.CalculateFeeRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CalculateFeeRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut amount__ = None;
                let mut payload_type__ = None;
                let mut fixed_amount__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PayloadType => {
                            if payload_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("payloadType"));
                            }
                            payload_type__ = Some(map_.next_value::<PayloadType>()? as i32);
                        }
                        GeneratedField::FixedAmount => {
                            if fixed_amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fixedAmount"));
                            }
                            fixed_amount__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(CalculateFeeRequest {
                    amount: amount__.unwrap_or_default(),
                    payload_type: payload_type__.unwrap_or_default(),
                    fixed_amount: fixed_amount__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.CalculateFeeRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CalculateFeeResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.amount != 0 {
            len += 1;
        }
        if self.fee != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.CalculateFeeResponse", len)?;
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        if self.fee != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("fee", ToString::to_string(&self.fee).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CalculateFeeResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "amount",
            "fee",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Amount,
            Fee,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "amount" => Ok(GeneratedField::Amount),
                            "fee" => Ok(GeneratedField::Fee),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CalculateFeeResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.CalculateFeeResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CalculateFeeResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut amount__ = None;
                let mut fee__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Fee => {
                            if fee__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fee"));
                            }
                            fee__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(CalculateFeeResponse {
                    amount: amount__.unwrap_or_default(),
                    fee: fee__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.CalculateFeeResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CertificateInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.hash.is_empty() {
            len += 1;
        }
        if self.round != 0 {
            len += 1;
        }
        if !self.committers.is_empty() {
            len += 1;
        }
        if !self.absentees.is_empty() {
            len += 1;
        }
        if !self.signature.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.CertificateInfo", len)?;
        if !self.hash.is_empty() {
            struct_ser.serialize_field("hash", &self.hash)?;
        }
        if self.round != 0 {
            struct_ser.serialize_field("round", &self.round)?;
        }
        if !self.committers.is_empty() {
            struct_ser.serialize_field("committers", &self.committers)?;
        }
        if !self.absentees.is_empty() {
            struct_ser.serialize_field("absentees", &self.absentees)?;
        }
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CertificateInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "hash",
            "round",
            "committers",
            "absentees",
            "signature",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Hash,
            Round,
            Committers,
            Absentees,
            Signature,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "hash" => Ok(GeneratedField::Hash),
                            "round" => Ok(GeneratedField::Round),
                            "committers" => Ok(GeneratedField::Committers),
                            "absentees" => Ok(GeneratedField::Absentees),
                            "signature" => Ok(GeneratedField::Signature),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CertificateInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.CertificateInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CertificateInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut hash__ = None;
                let mut round__ = None;
                let mut committers__ = None;
                let mut absentees__ = None;
                let mut signature__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Hash => {
                            if hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hash"));
                            }
                            hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Round => {
                            if round__.is_some() {
                                return Err(serde::de::Error::duplicate_field("round"));
                            }
                            round__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Committers => {
                            if committers__.is_some() {
                                return Err(serde::de::Error::duplicate_field("committers"));
                            }
                            committers__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::Absentees => {
                            if absentees__.is_some() {
                                return Err(serde::de::Error::duplicate_field("absentees"));
                            }
                            absentees__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(CertificateInfo {
                    hash: hash__.unwrap_or_default(),
                    round: round__.unwrap_or_default(),
                    committers: committers__.unwrap_or_default(),
                    absentees: absentees__.unwrap_or_default(),
                    signature: signature__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.CertificateInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ConnectionInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.connections != 0 {
            len += 1;
        }
        if self.inbound_connections != 0 {
            len += 1;
        }
        if self.outbound_connections != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ConnectionInfo", len)?;
        if self.connections != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("connections", ToString::to_string(&self.connections).as_str())?;
        }
        if self.inbound_connections != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("inboundConnections", ToString::to_string(&self.inbound_connections).as_str())?;
        }
        if self.outbound_connections != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("outboundConnections", ToString::to_string(&self.outbound_connections).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ConnectionInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "connections",
            "inbound_connections",
            "inboundConnections",
            "outbound_connections",
            "outboundConnections",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Connections,
            InboundConnections,
            OutboundConnections,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "connections" => Ok(GeneratedField::Connections),
                            "inboundConnections" | "inbound_connections" => Ok(GeneratedField::InboundConnections),
                            "outboundConnections" | "outbound_connections" => Ok(GeneratedField::OutboundConnections),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ConnectionInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ConnectionInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ConnectionInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut connections__ = None;
                let mut inbound_connections__ = None;
                let mut outbound_connections__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Connections => {
                            if connections__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connections"));
                            }
                            connections__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::InboundConnections => {
                            if inbound_connections__.is_some() {
                                return Err(serde::de::Error::duplicate_field("inboundConnections"));
                            }
                            inbound_connections__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::OutboundConnections => {
                            if outbound_connections__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outboundConnections"));
                            }
                            outbound_connections__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(ConnectionInfo {
                    connections: connections__.unwrap_or_default(),
                    inbound_connections: inbound_connections__.unwrap_or_default(),
                    outbound_connections: outbound_connections__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ConnectionInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ConsensusInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        if self.active {
            len += 1;
        }
        if self.height != 0 {
            len += 1;
        }
        if self.round != 0 {
            len += 1;
        }
        if !self.votes.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ConsensusInfo", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if self.active {
            struct_ser.serialize_field("active", &self.active)?;
        }
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        if self.round != 0 {
            struct_ser.serialize_field("round", &self.round)?;
        }
        if !self.votes.is_empty() {
            struct_ser.serialize_field("votes", &self.votes)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ConsensusInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
            "active",
            "height",
            "round",
            "votes",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
            Active,
            Height,
            Round,
            Votes,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            "active" => Ok(GeneratedField::Active),
                            "height" => Ok(GeneratedField::Height),
                            "round" => Ok(GeneratedField::Round),
                            "votes" => Ok(GeneratedField::Votes),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ConsensusInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ConsensusInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ConsensusInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                let mut active__ = None;
                let mut height__ = None;
                let mut round__ = None;
                let mut votes__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Active => {
                            if active__.is_some() {
                                return Err(serde::de::Error::duplicate_field("active"));
                            }
                            active__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Round => {
                            if round__.is_some() {
                                return Err(serde::de::Error::duplicate_field("round"));
                            }
                            round__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Votes => {
                            if votes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("votes"));
                            }
                            votes__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ConsensusInfo {
                    address: address__.unwrap_or_default(),
                    active: active__.unwrap_or_default(),
                    height: height__.unwrap_or_default(),
                    round: round__.unwrap_or_default(),
                    votes: votes__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ConsensusInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CounterInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.bytes != 0 {
            len += 1;
        }
        if self.bundles != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.CounterInfo", len)?;
        if self.bytes != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("bytes", ToString::to_string(&self.bytes).as_str())?;
        }
        if self.bundles != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("bundles", ToString::to_string(&self.bundles).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CounterInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "bytes",
            "bundles",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Bytes,
            Bundles,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "bytes" => Ok(GeneratedField::Bytes),
                            "bundles" => Ok(GeneratedField::Bundles),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CounterInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.CounterInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CounterInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut bytes__ = None;
                let mut bundles__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Bytes => {
                            if bytes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("bytes"));
                            }
                            bytes__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Bundles => {
                            if bundles__.is_some() {
                                return Err(serde::de::Error::duplicate_field("bundles"));
                            }
                            bundles__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(CounterInfo {
                    bytes: bytes__.unwrap_or_default(),
                    bundles: bundles__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.CounterInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateWalletRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.password.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.CreateWalletRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.password.is_empty() {
            struct_ser.serialize_field("password", &self.password)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateWalletRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "password",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Password,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "password" => Ok(GeneratedField::Password),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateWalletRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.CreateWalletRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateWalletRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut password__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Password => {
                            if password__.is_some() {
                                return Err(serde::de::Error::duplicate_field("password"));
                            }
                            password__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(CreateWalletRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    password: password__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.CreateWalletRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CreateWalletResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.mnemonic.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.CreateWalletResponse", len)?;
        if !self.mnemonic.is_empty() {
            struct_ser.serialize_field("mnemonic", &self.mnemonic)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CreateWalletResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "mnemonic",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Mnemonic,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "mnemonic" => Ok(GeneratedField::Mnemonic),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CreateWalletResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.CreateWalletResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CreateWalletResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut mnemonic__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Mnemonic => {
                            if mnemonic__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mnemonic"));
                            }
                            mnemonic__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(CreateWalletResponse {
                    mnemonic: mnemonic__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.CreateWalletResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for DecodeRawTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.raw_transaction.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.DecodeRawTransactionRequest", len)?;
        if !self.raw_transaction.is_empty() {
            struct_ser.serialize_field("rawTransaction", &self.raw_transaction)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for DecodeRawTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "raw_transaction",
            "rawTransaction",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RawTransaction,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "rawTransaction" | "raw_transaction" => Ok(GeneratedField::RawTransaction),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DecodeRawTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.DecodeRawTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<DecodeRawTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut raw_transaction__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RawTransaction => {
                            if raw_transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("rawTransaction"));
                            }
                            raw_transaction__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(DecodeRawTransactionRequest {
                    raw_transaction: raw_transaction__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.DecodeRawTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for DecodeRawTransactionResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.transaction.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.DecodeRawTransactionResponse", len)?;
        if let Some(v) = self.transaction.as_ref() {
            struct_ser.serialize_field("transaction", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for DecodeRawTransactionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "transaction",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Transaction,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "transaction" => Ok(GeneratedField::Transaction),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DecodeRawTransactionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.DecodeRawTransactionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<DecodeRawTransactionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut transaction__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Transaction => {
                            if transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("transaction"));
                            }
                            transaction__ = map_.next_value()?;
                        }
                    }
                }
                Ok(DecodeRawTransactionResponse {
                    transaction: transaction__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.DecodeRawTransactionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAccountRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetAccountRequest", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAccountRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAccountRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetAccountRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAccountRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAccountRequest {
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetAccountRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAccountResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.account.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetAccountResponse", len)?;
        if let Some(v) = self.account.as_ref() {
            struct_ser.serialize_field("account", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAccountResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "account",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Account,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "account" => Ok(GeneratedField::Account),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAccountResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetAccountResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAccountResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut account__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Account => {
                            if account__.is_some() {
                                return Err(serde::de::Error::duplicate_field("account"));
                            }
                            account__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetAccountResponse {
                    account: account__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetAccountResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAddressHistoryRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetAddressHistoryRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAddressHistoryRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAddressHistoryRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetAddressHistoryRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAddressHistoryRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAddressHistoryRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetAddressHistoryRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAddressHistoryResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.history_info.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetAddressHistoryResponse", len)?;
        if !self.history_info.is_empty() {
            struct_ser.serialize_field("historyInfo", &self.history_info)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAddressHistoryResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "history_info",
            "historyInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            HistoryInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "historyInfo" | "history_info" => Ok(GeneratedField::HistoryInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAddressHistoryResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetAddressHistoryResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAddressHistoryResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut history_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::HistoryInfo => {
                            if history_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("historyInfo"));
                            }
                            history_info__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAddressHistoryResponse {
                    history_info: history_info__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetAddressHistoryResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAddressInfoRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetAddressInfoRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAddressInfoRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAddressInfoRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetAddressInfoRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAddressInfoRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAddressInfoRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetAddressInfoRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAddressInfoResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if !self.label.is_empty() {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        if !self.path.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetAddressInfoResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if !self.label.is_empty() {
            struct_ser.serialize_field("label", &self.label)?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        if !self.path.is_empty() {
            struct_ser.serialize_field("path", &self.path)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAddressInfoResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "address",
            "label",
            "public_key",
            "publicKey",
            "path",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Address,
            Label,
            PublicKey,
            Path,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "address" => Ok(GeneratedField::Address),
                            "label" => Ok(GeneratedField::Label),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            "path" => Ok(GeneratedField::Path),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAddressInfoResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetAddressInfoResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAddressInfoResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut address__ = None;
                let mut label__ = None;
                let mut public_key__ = None;
                let mut path__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Label => {
                            if label__.is_some() {
                                return Err(serde::de::Error::duplicate_field("label"));
                            }
                            label__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Path => {
                            if path__.is_some() {
                                return Err(serde::de::Error::duplicate_field("path"));
                            }
                            path__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAddressInfoResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    label: label__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                    path: path__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetAddressInfoResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockHashRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.height != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockHashRequest", len)?;
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockHashRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "height",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Height,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "height" => Ok(GeneratedField::Height),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockHashRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockHashRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockHashRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut height__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetBlockHashRequest {
                    height: height__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockHashRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockHashResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.hash.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockHashResponse", len)?;
        if !self.hash.is_empty() {
            struct_ser.serialize_field("hash", &self.hash)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockHashResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "hash",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Hash,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "hash" => Ok(GeneratedField::Hash),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockHashResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockHashResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockHashResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut hash__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Hash => {
                            if hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hash"));
                            }
                            hash__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetBlockHashResponse {
                    hash: hash__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockHashResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockHeightRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.hash.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockHeightRequest", len)?;
        if !self.hash.is_empty() {
            struct_ser.serialize_field("hash", &self.hash)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockHeightRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "hash",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Hash,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "hash" => Ok(GeneratedField::Hash),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockHeightRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockHeightRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockHeightRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut hash__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Hash => {
                            if hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hash"));
                            }
                            hash__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetBlockHeightRequest {
                    hash: hash__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockHeightRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockHeightResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.height != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockHeightResponse", len)?;
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockHeightResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "height",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Height,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "height" => Ok(GeneratedField::Height),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockHeightResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockHeightResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockHeightResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut height__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetBlockHeightResponse {
                    height: height__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockHeightResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.height != 0 {
            len += 1;
        }
        if self.verbosity != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockRequest", len)?;
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        if self.verbosity != 0 {
            let v = BlockVerbosity::try_from(self.verbosity)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.verbosity)))?;
            struct_ser.serialize_field("verbosity", &v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "height",
            "verbosity",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Height,
            Verbosity,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "height" => Ok(GeneratedField::Height),
                            "verbosity" => Ok(GeneratedField::Verbosity),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut height__ = None;
                let mut verbosity__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Verbosity => {
                            if verbosity__.is_some() {
                                return Err(serde::de::Error::duplicate_field("verbosity"));
                            }
                            verbosity__ = Some(map_.next_value::<BlockVerbosity>()? as i32);
                        }
                    }
                }
                Ok(GetBlockRequest {
                    height: height__.unwrap_or_default(),
                    verbosity: verbosity__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.height != 0 {
            len += 1;
        }
        if !self.hash.is_empty() {
            len += 1;
        }
        if !self.data.is_empty() {
            len += 1;
        }
        if self.block_time != 0 {
            len += 1;
        }
        if self.header.is_some() {
            len += 1;
        }
        if self.prev_cert.is_some() {
            len += 1;
        }
        if !self.txs.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockResponse", len)?;
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        if !self.hash.is_empty() {
            struct_ser.serialize_field("hash", &self.hash)?;
        }
        if !self.data.is_empty() {
            struct_ser.serialize_field("data", &self.data)?;
        }
        if self.block_time != 0 {
            struct_ser.serialize_field("blockTime", &self.block_time)?;
        }
        if let Some(v) = self.header.as_ref() {
            struct_ser.serialize_field("header", v)?;
        }
        if let Some(v) = self.prev_cert.as_ref() {
            struct_ser.serialize_field("prevCert", v)?;
        }
        if !self.txs.is_empty() {
            struct_ser.serialize_field("txs", &self.txs)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "height",
            "hash",
            "data",
            "block_time",
            "blockTime",
            "header",
            "prev_cert",
            "prevCert",
            "txs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Height,
            Hash,
            Data,
            BlockTime,
            Header,
            PrevCert,
            Txs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "height" => Ok(GeneratedField::Height),
                            "hash" => Ok(GeneratedField::Hash),
                            "data" => Ok(GeneratedField::Data),
                            "blockTime" | "block_time" => Ok(GeneratedField::BlockTime),
                            "header" => Ok(GeneratedField::Header),
                            "prevCert" | "prev_cert" => Ok(GeneratedField::PrevCert),
                            "txs" => Ok(GeneratedField::Txs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut height__ = None;
                let mut hash__ = None;
                let mut data__ = None;
                let mut block_time__ = None;
                let mut header__ = None;
                let mut prev_cert__ = None;
                let mut txs__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Hash => {
                            if hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hash"));
                            }
                            hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Data => {
                            if data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("data"));
                            }
                            data__ = Some(map_.next_value()?);
                        }
                        GeneratedField::BlockTime => {
                            if block_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("blockTime"));
                            }
                            block_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Header => {
                            if header__.is_some() {
                                return Err(serde::de::Error::duplicate_field("header"));
                            }
                            header__ = map_.next_value()?;
                        }
                        GeneratedField::PrevCert => {
                            if prev_cert__.is_some() {
                                return Err(serde::de::Error::duplicate_field("prevCert"));
                            }
                            prev_cert__ = map_.next_value()?;
                        }
                        GeneratedField::Txs => {
                            if txs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txs"));
                            }
                            txs__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetBlockResponse {
                    height: height__.unwrap_or_default(),
                    hash: hash__.unwrap_or_default(),
                    data: data__.unwrap_or_default(),
                    block_time: block_time__.unwrap_or_default(),
                    header: header__,
                    prev_cert: prev_cert__,
                    txs: txs__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockchainInfoRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("pactus.GetBlockchainInfoRequest", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockchainInfoRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockchainInfoRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockchainInfoRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockchainInfoRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(GetBlockchainInfoRequest {
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockchainInfoRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetBlockchainInfoResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.last_block_height != 0 {
            len += 1;
        }
        if !self.last_block_hash.is_empty() {
            len += 1;
        }
        if self.total_accounts != 0 {
            len += 1;
        }
        if self.total_validators != 0 {
            len += 1;
        }
        if self.total_power != 0 {
            len += 1;
        }
        if self.committee_power != 0 {
            len += 1;
        }
        if !self.committee_validators.is_empty() {
            len += 1;
        }
        if self.is_pruned {
            len += 1;
        }
        if self.pruning_height != 0 {
            len += 1;
        }
        if self.last_block_time != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetBlockchainInfoResponse", len)?;
        if self.last_block_height != 0 {
            struct_ser.serialize_field("lastBlockHeight", &self.last_block_height)?;
        }
        if !self.last_block_hash.is_empty() {
            struct_ser.serialize_field("lastBlockHash", &self.last_block_hash)?;
        }
        if self.total_accounts != 0 {
            struct_ser.serialize_field("totalAccounts", &self.total_accounts)?;
        }
        if self.total_validators != 0 {
            struct_ser.serialize_field("totalValidators", &self.total_validators)?;
        }
        if self.total_power != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("totalPower", ToString::to_string(&self.total_power).as_str())?;
        }
        if self.committee_power != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("committeePower", ToString::to_string(&self.committee_power).as_str())?;
        }
        if !self.committee_validators.is_empty() {
            struct_ser.serialize_field("committeeValidators", &self.committee_validators)?;
        }
        if self.is_pruned {
            struct_ser.serialize_field("isPruned", &self.is_pruned)?;
        }
        if self.pruning_height != 0 {
            struct_ser.serialize_field("pruningHeight", &self.pruning_height)?;
        }
        if self.last_block_time != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("lastBlockTime", ToString::to_string(&self.last_block_time).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetBlockchainInfoResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "last_block_height",
            "lastBlockHeight",
            "last_block_hash",
            "lastBlockHash",
            "total_accounts",
            "totalAccounts",
            "total_validators",
            "totalValidators",
            "total_power",
            "totalPower",
            "committee_power",
            "committeePower",
            "committee_validators",
            "committeeValidators",
            "is_pruned",
            "isPruned",
            "pruning_height",
            "pruningHeight",
            "last_block_time",
            "lastBlockTime",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LastBlockHeight,
            LastBlockHash,
            TotalAccounts,
            TotalValidators,
            TotalPower,
            CommitteePower,
            CommitteeValidators,
            IsPruned,
            PruningHeight,
            LastBlockTime,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lastBlockHeight" | "last_block_height" => Ok(GeneratedField::LastBlockHeight),
                            "lastBlockHash" | "last_block_hash" => Ok(GeneratedField::LastBlockHash),
                            "totalAccounts" | "total_accounts" => Ok(GeneratedField::TotalAccounts),
                            "totalValidators" | "total_validators" => Ok(GeneratedField::TotalValidators),
                            "totalPower" | "total_power" => Ok(GeneratedField::TotalPower),
                            "committeePower" | "committee_power" => Ok(GeneratedField::CommitteePower),
                            "committeeValidators" | "committee_validators" => Ok(GeneratedField::CommitteeValidators),
                            "isPruned" | "is_pruned" => Ok(GeneratedField::IsPruned),
                            "pruningHeight" | "pruning_height" => Ok(GeneratedField::PruningHeight),
                            "lastBlockTime" | "last_block_time" => Ok(GeneratedField::LastBlockTime),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetBlockchainInfoResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetBlockchainInfoResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetBlockchainInfoResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut last_block_height__ = None;
                let mut last_block_hash__ = None;
                let mut total_accounts__ = None;
                let mut total_validators__ = None;
                let mut total_power__ = None;
                let mut committee_power__ = None;
                let mut committee_validators__ = None;
                let mut is_pruned__ = None;
                let mut pruning_height__ = None;
                let mut last_block_time__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LastBlockHeight => {
                            if last_block_height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastBlockHeight"));
                            }
                            last_block_height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastBlockHash => {
                            if last_block_hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastBlockHash"));
                            }
                            last_block_hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TotalAccounts => {
                            if total_accounts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalAccounts"));
                            }
                            total_accounts__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::TotalValidators => {
                            if total_validators__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalValidators"));
                            }
                            total_validators__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::TotalPower => {
                            if total_power__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalPower"));
                            }
                            total_power__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::CommitteePower => {
                            if committee_power__.is_some() {
                                return Err(serde::de::Error::duplicate_field("committeePower"));
                            }
                            committee_power__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::CommitteeValidators => {
                            if committee_validators__.is_some() {
                                return Err(serde::de::Error::duplicate_field("committeeValidators"));
                            }
                            committee_validators__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IsPruned => {
                            if is_pruned__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isPruned"));
                            }
                            is_pruned__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PruningHeight => {
                            if pruning_height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pruningHeight"));
                            }
                            pruning_height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastBlockTime => {
                            if last_block_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastBlockTime"));
                            }
                            last_block_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetBlockchainInfoResponse {
                    last_block_height: last_block_height__.unwrap_or_default(),
                    last_block_hash: last_block_hash__.unwrap_or_default(),
                    total_accounts: total_accounts__.unwrap_or_default(),
                    total_validators: total_validators__.unwrap_or_default(),
                    total_power: total_power__.unwrap_or_default(),
                    committee_power: committee_power__.unwrap_or_default(),
                    committee_validators: committee_validators__.unwrap_or_default(),
                    is_pruned: is_pruned__.unwrap_or_default(),
                    pruning_height: pruning_height__.unwrap_or_default(),
                    last_block_time: last_block_time__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetBlockchainInfoResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetConsensusInfoRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("pactus.GetConsensusInfoRequest", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetConsensusInfoRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetConsensusInfoRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetConsensusInfoRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetConsensusInfoRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(GetConsensusInfoRequest {
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetConsensusInfoRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetConsensusInfoResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.proposal.is_some() {
            len += 1;
        }
        if !self.instances.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetConsensusInfoResponse", len)?;
        if let Some(v) = self.proposal.as_ref() {
            struct_ser.serialize_field("proposal", v)?;
        }
        if !self.instances.is_empty() {
            struct_ser.serialize_field("instances", &self.instances)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetConsensusInfoResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "proposal",
            "instances",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Proposal,
            Instances,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "proposal" => Ok(GeneratedField::Proposal),
                            "instances" => Ok(GeneratedField::Instances),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetConsensusInfoResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetConsensusInfoResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetConsensusInfoResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut proposal__ = None;
                let mut instances__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Proposal => {
                            if proposal__.is_some() {
                                return Err(serde::de::Error::duplicate_field("proposal"));
                            }
                            proposal__ = map_.next_value()?;
                        }
                        GeneratedField::Instances => {
                            if instances__.is_some() {
                                return Err(serde::de::Error::duplicate_field("instances"));
                            }
                            instances__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetConsensusInfoResponse {
                    proposal: proposal__,
                    instances: instances__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetConsensusInfoResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetNetworkInfoRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.only_connected {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetNetworkInfoRequest", len)?;
        if self.only_connected {
            struct_ser.serialize_field("onlyConnected", &self.only_connected)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetNetworkInfoRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "only_connected",
            "onlyConnected",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            OnlyConnected,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "onlyConnected" | "only_connected" => Ok(GeneratedField::OnlyConnected),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetNetworkInfoRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetNetworkInfoRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetNetworkInfoRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut only_connected__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::OnlyConnected => {
                            if only_connected__.is_some() {
                                return Err(serde::de::Error::duplicate_field("onlyConnected"));
                            }
                            only_connected__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetNetworkInfoRequest {
                    only_connected: only_connected__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetNetworkInfoRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetNetworkInfoResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.network_name.is_empty() {
            len += 1;
        }
        if self.connected_peers_count != 0 {
            len += 1;
        }
        if !self.connected_peers.is_empty() {
            len += 1;
        }
        if self.metric_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetNetworkInfoResponse", len)?;
        if !self.network_name.is_empty() {
            struct_ser.serialize_field("networkName", &self.network_name)?;
        }
        if self.connected_peers_count != 0 {
            struct_ser.serialize_field("connectedPeersCount", &self.connected_peers_count)?;
        }
        if !self.connected_peers.is_empty() {
            struct_ser.serialize_field("connectedPeers", &self.connected_peers)?;
        }
        if let Some(v) = self.metric_info.as_ref() {
            struct_ser.serialize_field("metricInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetNetworkInfoResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "network_name",
            "networkName",
            "connected_peers_count",
            "connectedPeersCount",
            "connected_peers",
            "connectedPeers",
            "metric_info",
            "metricInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            NetworkName,
            ConnectedPeersCount,
            ConnectedPeers,
            MetricInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "networkName" | "network_name" => Ok(GeneratedField::NetworkName),
                            "connectedPeersCount" | "connected_peers_count" => Ok(GeneratedField::ConnectedPeersCount),
                            "connectedPeers" | "connected_peers" => Ok(GeneratedField::ConnectedPeers),
                            "metricInfo" | "metric_info" => Ok(GeneratedField::MetricInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetNetworkInfoResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetNetworkInfoResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetNetworkInfoResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut network_name__ = None;
                let mut connected_peers_count__ = None;
                let mut connected_peers__ = None;
                let mut metric_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::NetworkName => {
                            if network_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("networkName"));
                            }
                            network_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ConnectedPeersCount => {
                            if connected_peers_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connectedPeersCount"));
                            }
                            connected_peers_count__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ConnectedPeers => {
                            if connected_peers__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connectedPeers"));
                            }
                            connected_peers__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MetricInfo => {
                            if metric_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metricInfo"));
                            }
                            metric_info__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetNetworkInfoResponse {
                    network_name: network_name__.unwrap_or_default(),
                    connected_peers_count: connected_peers_count__.unwrap_or_default(),
                    connected_peers: connected_peers__.unwrap_or_default(),
                    metric_info: metric_info__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetNetworkInfoResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetNewAddressRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if self.address_type != 0 {
            len += 1;
        }
        if !self.label.is_empty() {
            len += 1;
        }
        if !self.password.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetNewAddressRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if self.address_type != 0 {
            let v = AddressType::try_from(self.address_type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.address_type)))?;
            struct_ser.serialize_field("addressType", &v)?;
        }
        if !self.label.is_empty() {
            struct_ser.serialize_field("label", &self.label)?;
        }
        if !self.password.is_empty() {
            struct_ser.serialize_field("password", &self.password)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetNewAddressRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "address_type",
            "addressType",
            "label",
            "password",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            AddressType,
            Label,
            Password,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "addressType" | "address_type" => Ok(GeneratedField::AddressType),
                            "label" => Ok(GeneratedField::Label),
                            "password" => Ok(GeneratedField::Password),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetNewAddressRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetNewAddressRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetNewAddressRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut address_type__ = None;
                let mut label__ = None;
                let mut password__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::AddressType => {
                            if address_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("addressType"));
                            }
                            address_type__ = Some(map_.next_value::<AddressType>()? as i32);
                        }
                        GeneratedField::Label => {
                            if label__.is_some() {
                                return Err(serde::de::Error::duplicate_field("label"));
                            }
                            label__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Password => {
                            if password__.is_some() {
                                return Err(serde::de::Error::duplicate_field("password"));
                            }
                            password__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetNewAddressRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    address_type: address_type__.unwrap_or_default(),
                    label: label__.unwrap_or_default(),
                    password: password__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetNewAddressRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetNewAddressResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if self.address_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetNewAddressResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if let Some(v) = self.address_info.as_ref() {
            struct_ser.serialize_field("addressInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetNewAddressResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "address_info",
            "addressInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            AddressInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "addressInfo" | "address_info" => Ok(GeneratedField::AddressInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetNewAddressResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetNewAddressResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetNewAddressResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut address_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::AddressInfo => {
                            if address_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("addressInfo"));
                            }
                            address_info__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetNewAddressResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    address_info: address_info__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetNewAddressResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetNodeInfoRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("pactus.GetNodeInfoRequest", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetNodeInfoRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetNodeInfoRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetNodeInfoRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetNodeInfoRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(GetNodeInfoRequest {
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetNodeInfoRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetNodeInfoResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.moniker.is_empty() {
            len += 1;
        }
        if !self.agent.is_empty() {
            len += 1;
        }
        if !self.peer_id.is_empty() {
            len += 1;
        }
        if self.started_at != 0 {
            len += 1;
        }
        if !self.reachability.is_empty() {
            len += 1;
        }
        if self.services != 0 {
            len += 1;
        }
        if !self.services_names.is_empty() {
            len += 1;
        }
        if !self.local_addrs.is_empty() {
            len += 1;
        }
        if !self.protocols.is_empty() {
            len += 1;
        }
        if self.clock_offset != 0. {
            len += 1;
        }
        if self.connection_info.is_some() {
            len += 1;
        }
        if !self.zmq_publishers.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetNodeInfoResponse", len)?;
        if !self.moniker.is_empty() {
            struct_ser.serialize_field("moniker", &self.moniker)?;
        }
        if !self.agent.is_empty() {
            struct_ser.serialize_field("agent", &self.agent)?;
        }
        if !self.peer_id.is_empty() {
            struct_ser.serialize_field("peerId", &self.peer_id)?;
        }
        if self.started_at != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startedAt", ToString::to_string(&self.started_at).as_str())?;
        }
        if !self.reachability.is_empty() {
            struct_ser.serialize_field("reachability", &self.reachability)?;
        }
        if self.services != 0 {
            struct_ser.serialize_field("services", &self.services)?;
        }
        if !self.services_names.is_empty() {
            struct_ser.serialize_field("servicesNames", &self.services_names)?;
        }
        if !self.local_addrs.is_empty() {
            struct_ser.serialize_field("localAddrs", &self.local_addrs)?;
        }
        if !self.protocols.is_empty() {
            struct_ser.serialize_field("protocols", &self.protocols)?;
        }
        if self.clock_offset != 0. {
            struct_ser.serialize_field("clockOffset", &self.clock_offset)?;
        }
        if let Some(v) = self.connection_info.as_ref() {
            struct_ser.serialize_field("connectionInfo", v)?;
        }
        if !self.zmq_publishers.is_empty() {
            struct_ser.serialize_field("zmqPublishers", &self.zmq_publishers)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetNodeInfoResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "moniker",
            "agent",
            "peer_id",
            "peerId",
            "started_at",
            "startedAt",
            "reachability",
            "services",
            "services_names",
            "servicesNames",
            "local_addrs",
            "localAddrs",
            "protocols",
            "clock_offset",
            "clockOffset",
            "connection_info",
            "connectionInfo",
            "zmq_publishers",
            "zmqPublishers",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Moniker,
            Agent,
            PeerId,
            StartedAt,
            Reachability,
            Services,
            ServicesNames,
            LocalAddrs,
            Protocols,
            ClockOffset,
            ConnectionInfo,
            ZmqPublishers,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "moniker" => Ok(GeneratedField::Moniker),
                            "agent" => Ok(GeneratedField::Agent),
                            "peerId" | "peer_id" => Ok(GeneratedField::PeerId),
                            "startedAt" | "started_at" => Ok(GeneratedField::StartedAt),
                            "reachability" => Ok(GeneratedField::Reachability),
                            "services" => Ok(GeneratedField::Services),
                            "servicesNames" | "services_names" => Ok(GeneratedField::ServicesNames),
                            "localAddrs" | "local_addrs" => Ok(GeneratedField::LocalAddrs),
                            "protocols" => Ok(GeneratedField::Protocols),
                            "clockOffset" | "clock_offset" => Ok(GeneratedField::ClockOffset),
                            "connectionInfo" | "connection_info" => Ok(GeneratedField::ConnectionInfo),
                            "zmqPublishers" | "zmq_publishers" => Ok(GeneratedField::ZmqPublishers),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetNodeInfoResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetNodeInfoResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetNodeInfoResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut moniker__ = None;
                let mut agent__ = None;
                let mut peer_id__ = None;
                let mut started_at__ = None;
                let mut reachability__ = None;
                let mut services__ = None;
                let mut services_names__ = None;
                let mut local_addrs__ = None;
                let mut protocols__ = None;
                let mut clock_offset__ = None;
                let mut connection_info__ = None;
                let mut zmq_publishers__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Moniker => {
                            if moniker__.is_some() {
                                return Err(serde::de::Error::duplicate_field("moniker"));
                            }
                            moniker__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Agent => {
                            if agent__.is_some() {
                                return Err(serde::de::Error::duplicate_field("agent"));
                            }
                            agent__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PeerId => {
                            if peer_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("peerId"));
                            }
                            peer_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::StartedAt => {
                            if started_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startedAt"));
                            }
                            started_at__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Reachability => {
                            if reachability__.is_some() {
                                return Err(serde::de::Error::duplicate_field("reachability"));
                            }
                            reachability__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Services => {
                            if services__.is_some() {
                                return Err(serde::de::Error::duplicate_field("services"));
                            }
                            services__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ServicesNames => {
                            if services_names__.is_some() {
                                return Err(serde::de::Error::duplicate_field("servicesNames"));
                            }
                            services_names__ = Some(map_.next_value()?);
                        }
                        GeneratedField::LocalAddrs => {
                            if local_addrs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("localAddrs"));
                            }
                            local_addrs__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Protocols => {
                            if protocols__.is_some() {
                                return Err(serde::de::Error::duplicate_field("protocols"));
                            }
                            protocols__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ClockOffset => {
                            if clock_offset__.is_some() {
                                return Err(serde::de::Error::duplicate_field("clockOffset"));
                            }
                            clock_offset__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ConnectionInfo => {
                            if connection_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connectionInfo"));
                            }
                            connection_info__ = map_.next_value()?;
                        }
                        GeneratedField::ZmqPublishers => {
                            if zmq_publishers__.is_some() {
                                return Err(serde::de::Error::duplicate_field("zmqPublishers"));
                            }
                            zmq_publishers__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetNodeInfoResponse {
                    moniker: moniker__.unwrap_or_default(),
                    agent: agent__.unwrap_or_default(),
                    peer_id: peer_id__.unwrap_or_default(),
                    started_at: started_at__.unwrap_or_default(),
                    reachability: reachability__.unwrap_or_default(),
                    services: services__.unwrap_or_default(),
                    services_names: services_names__.unwrap_or_default(),
                    local_addrs: local_addrs__.unwrap_or_default(),
                    protocols: protocols__.unwrap_or_default(),
                    clock_offset: clock_offset__.unwrap_or_default(),
                    connection_info: connection_info__,
                    zmq_publishers: zmq_publishers__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetNodeInfoResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetPublicKeyRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetPublicKeyRequest", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetPublicKeyRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetPublicKeyRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetPublicKeyRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetPublicKeyRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetPublicKeyRequest {
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetPublicKeyRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetPublicKeyResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.public_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetPublicKeyResponse", len)?;
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetPublicKeyResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "public_key",
            "publicKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PublicKey,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetPublicKeyResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetPublicKeyResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetPublicKeyResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut public_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetPublicKeyResponse {
                    public_key: public_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetPublicKeyResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRawBondTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lock_time != 0 {
            len += 1;
        }
        if !self.sender.is_empty() {
            len += 1;
        }
        if !self.receiver.is_empty() {
            len += 1;
        }
        if self.stake != 0 {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        if self.fee != 0 {
            len += 1;
        }
        if !self.memo.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetRawBondTransactionRequest", len)?;
        if self.lock_time != 0 {
            struct_ser.serialize_field("lockTime", &self.lock_time)?;
        }
        if !self.sender.is_empty() {
            struct_ser.serialize_field("sender", &self.sender)?;
        }
        if !self.receiver.is_empty() {
            struct_ser.serialize_field("receiver", &self.receiver)?;
        }
        if self.stake != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("stake", ToString::to_string(&self.stake).as_str())?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        if self.fee != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("fee", ToString::to_string(&self.fee).as_str())?;
        }
        if !self.memo.is_empty() {
            struct_ser.serialize_field("memo", &self.memo)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRawBondTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "lock_time",
            "lockTime",
            "sender",
            "receiver",
            "stake",
            "public_key",
            "publicKey",
            "fee",
            "memo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LockTime,
            Sender,
            Receiver,
            Stake,
            PublicKey,
            Fee,
            Memo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lockTime" | "lock_time" => Ok(GeneratedField::LockTime),
                            "sender" => Ok(GeneratedField::Sender),
                            "receiver" => Ok(GeneratedField::Receiver),
                            "stake" => Ok(GeneratedField::Stake),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            "fee" => Ok(GeneratedField::Fee),
                            "memo" => Ok(GeneratedField::Memo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRawBondTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetRawBondTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRawBondTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut lock_time__ = None;
                let mut sender__ = None;
                let mut receiver__ = None;
                let mut stake__ = None;
                let mut public_key__ = None;
                let mut fee__ = None;
                let mut memo__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LockTime => {
                            if lock_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockTime"));
                            }
                            lock_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Sender => {
                            if sender__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sender"));
                            }
                            sender__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Receiver => {
                            if receiver__.is_some() {
                                return Err(serde::de::Error::duplicate_field("receiver"));
                            }
                            receiver__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Stake => {
                            if stake__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stake"));
                            }
                            stake__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Fee => {
                            if fee__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fee"));
                            }
                            fee__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Memo => {
                            if memo__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memo"));
                            }
                            memo__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetRawBondTransactionRequest {
                    lock_time: lock_time__.unwrap_or_default(),
                    sender: sender__.unwrap_or_default(),
                    receiver: receiver__.unwrap_or_default(),
                    stake: stake__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                    fee: fee__.unwrap_or_default(),
                    memo: memo__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetRawBondTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRawTransactionResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.raw_transaction.is_empty() {
            len += 1;
        }
        if !self.id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetRawTransactionResponse", len)?;
        if !self.raw_transaction.is_empty() {
            struct_ser.serialize_field("rawTransaction", &self.raw_transaction)?;
        }
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRawTransactionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "raw_transaction",
            "rawTransaction",
            "id",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RawTransaction,
            Id,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "rawTransaction" | "raw_transaction" => Ok(GeneratedField::RawTransaction),
                            "id" => Ok(GeneratedField::Id),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRawTransactionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetRawTransactionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRawTransactionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut raw_transaction__ = None;
                let mut id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RawTransaction => {
                            if raw_transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("rawTransaction"));
                            }
                            raw_transaction__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetRawTransactionResponse {
                    raw_transaction: raw_transaction__.unwrap_or_default(),
                    id: id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetRawTransactionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRawTransferTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lock_time != 0 {
            len += 1;
        }
        if !self.sender.is_empty() {
            len += 1;
        }
        if !self.receiver.is_empty() {
            len += 1;
        }
        if self.amount != 0 {
            len += 1;
        }
        if self.fee != 0 {
            len += 1;
        }
        if !self.memo.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetRawTransferTransactionRequest", len)?;
        if self.lock_time != 0 {
            struct_ser.serialize_field("lockTime", &self.lock_time)?;
        }
        if !self.sender.is_empty() {
            struct_ser.serialize_field("sender", &self.sender)?;
        }
        if !self.receiver.is_empty() {
            struct_ser.serialize_field("receiver", &self.receiver)?;
        }
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        if self.fee != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("fee", ToString::to_string(&self.fee).as_str())?;
        }
        if !self.memo.is_empty() {
            struct_ser.serialize_field("memo", &self.memo)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRawTransferTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "lock_time",
            "lockTime",
            "sender",
            "receiver",
            "amount",
            "fee",
            "memo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LockTime,
            Sender,
            Receiver,
            Amount,
            Fee,
            Memo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lockTime" | "lock_time" => Ok(GeneratedField::LockTime),
                            "sender" => Ok(GeneratedField::Sender),
                            "receiver" => Ok(GeneratedField::Receiver),
                            "amount" => Ok(GeneratedField::Amount),
                            "fee" => Ok(GeneratedField::Fee),
                            "memo" => Ok(GeneratedField::Memo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRawTransferTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetRawTransferTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRawTransferTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut lock_time__ = None;
                let mut sender__ = None;
                let mut receiver__ = None;
                let mut amount__ = None;
                let mut fee__ = None;
                let mut memo__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LockTime => {
                            if lock_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockTime"));
                            }
                            lock_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Sender => {
                            if sender__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sender"));
                            }
                            sender__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Receiver => {
                            if receiver__.is_some() {
                                return Err(serde::de::Error::duplicate_field("receiver"));
                            }
                            receiver__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Fee => {
                            if fee__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fee"));
                            }
                            fee__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Memo => {
                            if memo__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memo"));
                            }
                            memo__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetRawTransferTransactionRequest {
                    lock_time: lock_time__.unwrap_or_default(),
                    sender: sender__.unwrap_or_default(),
                    receiver: receiver__.unwrap_or_default(),
                    amount: amount__.unwrap_or_default(),
                    fee: fee__.unwrap_or_default(),
                    memo: memo__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetRawTransferTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRawUnbondTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lock_time != 0 {
            len += 1;
        }
        if !self.validator_address.is_empty() {
            len += 1;
        }
        if !self.memo.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetRawUnbondTransactionRequest", len)?;
        if self.lock_time != 0 {
            struct_ser.serialize_field("lockTime", &self.lock_time)?;
        }
        if !self.validator_address.is_empty() {
            struct_ser.serialize_field("validatorAddress", &self.validator_address)?;
        }
        if !self.memo.is_empty() {
            struct_ser.serialize_field("memo", &self.memo)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRawUnbondTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "lock_time",
            "lockTime",
            "validator_address",
            "validatorAddress",
            "memo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LockTime,
            ValidatorAddress,
            Memo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lockTime" | "lock_time" => Ok(GeneratedField::LockTime),
                            "validatorAddress" | "validator_address" => Ok(GeneratedField::ValidatorAddress),
                            "memo" => Ok(GeneratedField::Memo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRawUnbondTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetRawUnbondTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRawUnbondTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut lock_time__ = None;
                let mut validator_address__ = None;
                let mut memo__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LockTime => {
                            if lock_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockTime"));
                            }
                            lock_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ValidatorAddress => {
                            if validator_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("validatorAddress"));
                            }
                            validator_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Memo => {
                            if memo__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memo"));
                            }
                            memo__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetRawUnbondTransactionRequest {
                    lock_time: lock_time__.unwrap_or_default(),
                    validator_address: validator_address__.unwrap_or_default(),
                    memo: memo__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetRawUnbondTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetRawWithdrawTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lock_time != 0 {
            len += 1;
        }
        if !self.validator_address.is_empty() {
            len += 1;
        }
        if !self.account_address.is_empty() {
            len += 1;
        }
        if self.amount != 0 {
            len += 1;
        }
        if self.fee != 0 {
            len += 1;
        }
        if !self.memo.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetRawWithdrawTransactionRequest", len)?;
        if self.lock_time != 0 {
            struct_ser.serialize_field("lockTime", &self.lock_time)?;
        }
        if !self.validator_address.is_empty() {
            struct_ser.serialize_field("validatorAddress", &self.validator_address)?;
        }
        if !self.account_address.is_empty() {
            struct_ser.serialize_field("accountAddress", &self.account_address)?;
        }
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        if self.fee != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("fee", ToString::to_string(&self.fee).as_str())?;
        }
        if !self.memo.is_empty() {
            struct_ser.serialize_field("memo", &self.memo)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetRawWithdrawTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "lock_time",
            "lockTime",
            "validator_address",
            "validatorAddress",
            "account_address",
            "accountAddress",
            "amount",
            "fee",
            "memo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LockTime,
            ValidatorAddress,
            AccountAddress,
            Amount,
            Fee,
            Memo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lockTime" | "lock_time" => Ok(GeneratedField::LockTime),
                            "validatorAddress" | "validator_address" => Ok(GeneratedField::ValidatorAddress),
                            "accountAddress" | "account_address" => Ok(GeneratedField::AccountAddress),
                            "amount" => Ok(GeneratedField::Amount),
                            "fee" => Ok(GeneratedField::Fee),
                            "memo" => Ok(GeneratedField::Memo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetRawWithdrawTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetRawWithdrawTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetRawWithdrawTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut lock_time__ = None;
                let mut validator_address__ = None;
                let mut account_address__ = None;
                let mut amount__ = None;
                let mut fee__ = None;
                let mut memo__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LockTime => {
                            if lock_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockTime"));
                            }
                            lock_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ValidatorAddress => {
                            if validator_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("validatorAddress"));
                            }
                            validator_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::AccountAddress => {
                            if account_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("accountAddress"));
                            }
                            account_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Fee => {
                            if fee__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fee"));
                            }
                            fee__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Memo => {
                            if memo__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memo"));
                            }
                            memo__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetRawWithdrawTransactionRequest {
                    lock_time: lock_time__.unwrap_or_default(),
                    validator_address: validator_address__.unwrap_or_default(),
                    account_address: account_address__.unwrap_or_default(),
                    amount: amount__.unwrap_or_default(),
                    fee: fee__.unwrap_or_default(),
                    memo: memo__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetRawWithdrawTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTotalBalanceRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTotalBalanceRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTotalBalanceRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTotalBalanceRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTotalBalanceRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTotalBalanceRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetTotalBalanceRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTotalBalanceRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTotalBalanceResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if self.total_balance != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTotalBalanceResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if self.total_balance != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("totalBalance", ToString::to_string(&self.total_balance).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTotalBalanceResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "total_balance",
            "totalBalance",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            TotalBalance,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "totalBalance" | "total_balance" => Ok(GeneratedField::TotalBalance),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTotalBalanceResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTotalBalanceResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTotalBalanceResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut total_balance__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TotalBalance => {
                            if total_balance__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalBalance"));
                            }
                            total_balance__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetTotalBalanceResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    total_balance: total_balance__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTotalBalanceResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTotalStakeRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTotalStakeRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTotalStakeRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTotalStakeRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTotalStakeRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTotalStakeRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetTotalStakeRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTotalStakeRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTotalStakeResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if self.total_stake != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTotalStakeResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if self.total_stake != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("totalStake", ToString::to_string(&self.total_stake).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTotalStakeResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "total_stake",
            "totalStake",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            TotalStake,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "totalStake" | "total_stake" => Ok(GeneratedField::TotalStake),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTotalStakeResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTotalStakeResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTotalStakeResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut total_stake__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TotalStake => {
                            if total_stake__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalStake"));
                            }
                            total_stake__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetTotalStakeResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    total_stake: total_stake__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTotalStakeResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if self.verbosity != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTransactionRequest", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if self.verbosity != 0 {
            let v = TransactionVerbosity::try_from(self.verbosity)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.verbosity)))?;
            struct_ser.serialize_field("verbosity", &v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "verbosity",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Verbosity,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "id" => Ok(GeneratedField::Id),
                            "verbosity" => Ok(GeneratedField::Verbosity),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut verbosity__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Verbosity => {
                            if verbosity__.is_some() {
                                return Err(serde::de::Error::duplicate_field("verbosity"));
                            }
                            verbosity__ = Some(map_.next_value::<TransactionVerbosity>()? as i32);
                        }
                    }
                }
                Ok(GetTransactionRequest {
                    id: id__.unwrap_or_default(),
                    verbosity: verbosity__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTransactionResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.block_height != 0 {
            len += 1;
        }
        if self.block_time != 0 {
            len += 1;
        }
        if self.transaction.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTransactionResponse", len)?;
        if self.block_height != 0 {
            struct_ser.serialize_field("blockHeight", &self.block_height)?;
        }
        if self.block_time != 0 {
            struct_ser.serialize_field("blockTime", &self.block_time)?;
        }
        if let Some(v) = self.transaction.as_ref() {
            struct_ser.serialize_field("transaction", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTransactionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "block_height",
            "blockHeight",
            "block_time",
            "blockTime",
            "transaction",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            BlockHeight,
            BlockTime,
            Transaction,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "blockHeight" | "block_height" => Ok(GeneratedField::BlockHeight),
                            "blockTime" | "block_time" => Ok(GeneratedField::BlockTime),
                            "transaction" => Ok(GeneratedField::Transaction),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTransactionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTransactionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTransactionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut block_height__ = None;
                let mut block_time__ = None;
                let mut transaction__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::BlockHeight => {
                            if block_height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("blockHeight"));
                            }
                            block_height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::BlockTime => {
                            if block_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("blockTime"));
                            }
                            block_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Transaction => {
                            if transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("transaction"));
                            }
                            transaction__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetTransactionResponse {
                    block_height: block_height__.unwrap_or_default(),
                    block_time: block_time__.unwrap_or_default(),
                    transaction: transaction__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTransactionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTxPoolContentRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload_type != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTxPoolContentRequest", len)?;
        if self.payload_type != 0 {
            let v = PayloadType::try_from(self.payload_type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.payload_type)))?;
            struct_ser.serialize_field("payloadType", &v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTxPoolContentRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "payload_type",
            "payloadType",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PayloadType,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "payloadType" | "payload_type" => Ok(GeneratedField::PayloadType),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTxPoolContentRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTxPoolContentRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTxPoolContentRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload_type__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PayloadType => {
                            if payload_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("payloadType"));
                            }
                            payload_type__ = Some(map_.next_value::<PayloadType>()? as i32);
                        }
                    }
                }
                Ok(GetTxPoolContentRequest {
                    payload_type: payload_type__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTxPoolContentRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetTxPoolContentResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.txs.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetTxPoolContentResponse", len)?;
        if !self.txs.is_empty() {
            struct_ser.serialize_field("txs", &self.txs)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetTxPoolContentResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "txs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Txs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "txs" => Ok(GeneratedField::Txs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetTxPoolContentResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetTxPoolContentResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetTxPoolContentResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut txs__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Txs => {
                            if txs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txs"));
                            }
                            txs__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetTxPoolContentResponse {
                    txs: txs__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetTxPoolContentResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorAddressRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.public_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetValidatorAddressRequest", len)?;
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorAddressRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "public_key",
            "publicKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PublicKey,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorAddressRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorAddressRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorAddressRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut public_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetValidatorAddressRequest {
                    public_key: public_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorAddressRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorAddressResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetValidatorAddressResponse", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorAddressResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorAddressResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorAddressResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorAddressResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetValidatorAddressResponse {
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorAddressResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorAddressesRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("pactus.GetValidatorAddressesRequest", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorAddressesRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorAddressesRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorAddressesRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorAddressesRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(GetValidatorAddressesRequest {
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorAddressesRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorAddressesResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.addresses.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetValidatorAddressesResponse", len)?;
        if !self.addresses.is_empty() {
            struct_ser.serialize_field("addresses", &self.addresses)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorAddressesResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "addresses",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Addresses,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "addresses" => Ok(GeneratedField::Addresses),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorAddressesResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorAddressesResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorAddressesResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut addresses__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Addresses => {
                            if addresses__.is_some() {
                                return Err(serde::de::Error::duplicate_field("addresses"));
                            }
                            addresses__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetValidatorAddressesResponse {
                    addresses: addresses__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorAddressesResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorByNumberRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.number != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetValidatorByNumberRequest", len)?;
        if self.number != 0 {
            struct_ser.serialize_field("number", &self.number)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorByNumberRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "number",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Number,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "number" => Ok(GeneratedField::Number),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorByNumberRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorByNumberRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorByNumberRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut number__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetValidatorByNumberRequest {
                    number: number__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorByNumberRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetValidatorRequest", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetValidatorRequest {
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetValidatorResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.validator.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetValidatorResponse", len)?;
        if let Some(v) = self.validator.as_ref() {
            struct_ser.serialize_field("validator", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetValidatorResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "validator",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Validator,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "validator" => Ok(GeneratedField::Validator),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetValidatorResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetValidatorResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetValidatorResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut validator__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Validator => {
                            if validator__.is_some() {
                                return Err(serde::de::Error::duplicate_field("validator"));
                            }
                            validator__ = map_.next_value()?;
                        }
                    }
                }
                Ok(GetValidatorResponse {
                    validator: validator__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetValidatorResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetWalletInfoRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetWalletInfoRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetWalletInfoRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetWalletInfoRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetWalletInfoRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetWalletInfoRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetWalletInfoRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetWalletInfoRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetWalletInfoResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if self.version != 0 {
            len += 1;
        }
        if !self.network.is_empty() {
            len += 1;
        }
        if self.encrypted {
            len += 1;
        }
        if !self.uuid.is_empty() {
            len += 1;
        }
        if self.created_at != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.GetWalletInfoResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if self.version != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("version", ToString::to_string(&self.version).as_str())?;
        }
        if !self.network.is_empty() {
            struct_ser.serialize_field("network", &self.network)?;
        }
        if self.encrypted {
            struct_ser.serialize_field("encrypted", &self.encrypted)?;
        }
        if !self.uuid.is_empty() {
            struct_ser.serialize_field("uuid", &self.uuid)?;
        }
        if self.created_at != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("createdAt", ToString::to_string(&self.created_at).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetWalletInfoResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "version",
            "network",
            "encrypted",
            "uuid",
            "created_at",
            "createdAt",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Version,
            Network,
            Encrypted,
            Uuid,
            CreatedAt,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "version" => Ok(GeneratedField::Version),
                            "network" => Ok(GeneratedField::Network),
                            "encrypted" => Ok(GeneratedField::Encrypted),
                            "uuid" => Ok(GeneratedField::Uuid),
                            "createdAt" | "created_at" => Ok(GeneratedField::CreatedAt),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetWalletInfoResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.GetWalletInfoResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetWalletInfoResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut version__ = None;
                let mut network__ = None;
                let mut encrypted__ = None;
                let mut uuid__ = None;
                let mut created_at__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Network => {
                            if network__.is_some() {
                                return Err(serde::de::Error::duplicate_field("network"));
                            }
                            network__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Encrypted => {
                            if encrypted__.is_some() {
                                return Err(serde::de::Error::duplicate_field("encrypted"));
                            }
                            encrypted__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Uuid => {
                            if uuid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("uuid"));
                            }
                            uuid__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CreatedAt => {
                            if created_at__.is_some() {
                                return Err(serde::de::Error::duplicate_field("createdAt"));
                            }
                            created_at__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(GetWalletInfoResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    version: version__.unwrap_or_default(),
                    network: network__.unwrap_or_default(),
                    encrypted: encrypted__.unwrap_or_default(),
                    uuid: uuid__.unwrap_or_default(),
                    created_at: created_at__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.GetWalletInfoResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for HistoryInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.transaction_id.is_empty() {
            len += 1;
        }
        if self.time != 0 {
            len += 1;
        }
        if !self.payload_type.is_empty() {
            len += 1;
        }
        if !self.description.is_empty() {
            len += 1;
        }
        if self.amount != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.HistoryInfo", len)?;
        if !self.transaction_id.is_empty() {
            struct_ser.serialize_field("transactionId", &self.transaction_id)?;
        }
        if self.time != 0 {
            struct_ser.serialize_field("time", &self.time)?;
        }
        if !self.payload_type.is_empty() {
            struct_ser.serialize_field("payloadType", &self.payload_type)?;
        }
        if !self.description.is_empty() {
            struct_ser.serialize_field("description", &self.description)?;
        }
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for HistoryInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "transaction_id",
            "transactionId",
            "time",
            "payload_type",
            "payloadType",
            "description",
            "amount",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TransactionId,
            Time,
            PayloadType,
            Description,
            Amount,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "transactionId" | "transaction_id" => Ok(GeneratedField::TransactionId),
                            "time" => Ok(GeneratedField::Time),
                            "payloadType" | "payload_type" => Ok(GeneratedField::PayloadType),
                            "description" => Ok(GeneratedField::Description),
                            "amount" => Ok(GeneratedField::Amount),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = HistoryInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.HistoryInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<HistoryInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut transaction_id__ = None;
                let mut time__ = None;
                let mut payload_type__ = None;
                let mut description__ = None;
                let mut amount__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TransactionId => {
                            if transaction_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("transactionId"));
                            }
                            transaction_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Time => {
                            if time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("time"));
                            }
                            time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PayloadType => {
                            if payload_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("payloadType"));
                            }
                            payload_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Description => {
                            if description__.is_some() {
                                return Err(serde::de::Error::duplicate_field("description"));
                            }
                            description__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(HistoryInfo {
                    transaction_id: transaction_id__.unwrap_or_default(),
                    time: time__.unwrap_or_default(),
                    payload_type: payload_type__.unwrap_or_default(),
                    description: description__.unwrap_or_default(),
                    amount: amount__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.HistoryInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListAddressRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ListAddressRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListAddressRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListAddressRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ListAddressRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListAddressRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ListAddressRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ListAddressRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListAddressResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.data.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ListAddressResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.data.is_empty() {
            struct_ser.serialize_field("data", &self.data)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListAddressResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "data",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Data,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "data" => Ok(GeneratedField::Data),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListAddressResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ListAddressResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListAddressResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut data__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Data => {
                            if data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("data"));
                            }
                            data__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ListAddressResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    data: data__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ListAddressResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListWalletRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("pactus.ListWalletRequest", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListWalletRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListWalletRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ListWalletRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListWalletRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(ListWalletRequest {
                })
            }
        }
        deserializer.deserialize_struct("pactus.ListWalletRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ListWalletResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallets.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ListWalletResponse", len)?;
        if !self.wallets.is_empty() {
            struct_ser.serialize_field("wallets", &self.wallets)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ListWalletResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallets",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Wallets,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "wallets" => Ok(GeneratedField::Wallets),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ListWalletResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ListWalletResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ListWalletResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallets__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Wallets => {
                            if wallets__.is_some() {
                                return Err(serde::de::Error::duplicate_field("wallets"));
                            }
                            wallets__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ListWalletResponse {
                    wallets: wallets__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ListWalletResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for LoadWalletRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.LoadWalletRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for LoadWalletRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = LoadWalletRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.LoadWalletRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<LoadWalletRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(LoadWalletRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.LoadWalletRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for LoadWalletResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.LoadWalletResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for LoadWalletResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = LoadWalletResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.LoadWalletResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<LoadWalletResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(LoadWalletResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.LoadWalletResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MetricInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.total_invalid.is_some() {
            len += 1;
        }
        if self.total_sent.is_some() {
            len += 1;
        }
        if self.total_received.is_some() {
            len += 1;
        }
        if !self.message_sent.is_empty() {
            len += 1;
        }
        if !self.message_received.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.MetricInfo", len)?;
        if let Some(v) = self.total_invalid.as_ref() {
            struct_ser.serialize_field("totalInvalid", v)?;
        }
        if let Some(v) = self.total_sent.as_ref() {
            struct_ser.serialize_field("totalSent", v)?;
        }
        if let Some(v) = self.total_received.as_ref() {
            struct_ser.serialize_field("totalReceived", v)?;
        }
        if !self.message_sent.is_empty() {
            struct_ser.serialize_field("messageSent", &self.message_sent)?;
        }
        if !self.message_received.is_empty() {
            struct_ser.serialize_field("messageReceived", &self.message_received)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MetricInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "total_invalid",
            "totalInvalid",
            "total_sent",
            "totalSent",
            "total_received",
            "totalReceived",
            "message_sent",
            "messageSent",
            "message_received",
            "messageReceived",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TotalInvalid,
            TotalSent,
            TotalReceived,
            MessageSent,
            MessageReceived,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "totalInvalid" | "total_invalid" => Ok(GeneratedField::TotalInvalid),
                            "totalSent" | "total_sent" => Ok(GeneratedField::TotalSent),
                            "totalReceived" | "total_received" => Ok(GeneratedField::TotalReceived),
                            "messageSent" | "message_sent" => Ok(GeneratedField::MessageSent),
                            "messageReceived" | "message_received" => Ok(GeneratedField::MessageReceived),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MetricInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.MetricInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MetricInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut total_invalid__ = None;
                let mut total_sent__ = None;
                let mut total_received__ = None;
                let mut message_sent__ = None;
                let mut message_received__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TotalInvalid => {
                            if total_invalid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalInvalid"));
                            }
                            total_invalid__ = map_.next_value()?;
                        }
                        GeneratedField::TotalSent => {
                            if total_sent__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalSent"));
                            }
                            total_sent__ = map_.next_value()?;
                        }
                        GeneratedField::TotalReceived => {
                            if total_received__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalReceived"));
                            }
                            total_received__ = map_.next_value()?;
                        }
                        GeneratedField::MessageSent => {
                            if message_sent__.is_some() {
                                return Err(serde::de::Error::duplicate_field("messageSent"));
                            }
                            message_sent__ = Some(
                                map_.next_value::<std::collections::HashMap<::pbjson::private::NumberDeserialize<i32>, _>>()?
                                    .into_iter().map(|(k,v)| (k.0, v)).collect()
                            );
                        }
                        GeneratedField::MessageReceived => {
                            if message_received__.is_some() {
                                return Err(serde::de::Error::duplicate_field("messageReceived"));
                            }
                            message_received__ = Some(
                                map_.next_value::<std::collections::HashMap<::pbjson::private::NumberDeserialize<i32>, _>>()?
                                    .into_iter().map(|(k,v)| (k.0, v)).collect()
                            );
                        }
                    }
                }
                Ok(MetricInfo {
                    total_invalid: total_invalid__,
                    total_sent: total_sent__,
                    total_received: total_received__,
                    message_sent: message_sent__.unwrap_or_default(),
                    message_received: message_received__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.MetricInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PayloadBond {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sender.is_empty() {
            len += 1;
        }
        if !self.receiver.is_empty() {
            len += 1;
        }
        if self.stake != 0 {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PayloadBond", len)?;
        if !self.sender.is_empty() {
            struct_ser.serialize_field("sender", &self.sender)?;
        }
        if !self.receiver.is_empty() {
            struct_ser.serialize_field("receiver", &self.receiver)?;
        }
        if self.stake != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("stake", ToString::to_string(&self.stake).as_str())?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PayloadBond {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sender",
            "receiver",
            "stake",
            "public_key",
            "publicKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sender,
            Receiver,
            Stake,
            PublicKey,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "sender" => Ok(GeneratedField::Sender),
                            "receiver" => Ok(GeneratedField::Receiver),
                            "stake" => Ok(GeneratedField::Stake),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PayloadBond;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PayloadBond")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PayloadBond, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sender__ = None;
                let mut receiver__ = None;
                let mut stake__ = None;
                let mut public_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sender => {
                            if sender__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sender"));
                            }
                            sender__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Receiver => {
                            if receiver__.is_some() {
                                return Err(serde::de::Error::duplicate_field("receiver"));
                            }
                            receiver__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Stake => {
                            if stake__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stake"));
                            }
                            stake__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(PayloadBond {
                    sender: sender__.unwrap_or_default(),
                    receiver: receiver__.unwrap_or_default(),
                    stake: stake__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PayloadBond", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PayloadSortition {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.address.is_empty() {
            len += 1;
        }
        if !self.proof.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PayloadSortition", len)?;
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if !self.proof.is_empty() {
            struct_ser.serialize_field("proof", &self.proof)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PayloadSortition {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "address",
            "proof",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Address,
            Proof,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "address" => Ok(GeneratedField::Address),
                            "proof" => Ok(GeneratedField::Proof),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PayloadSortition;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PayloadSortition")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PayloadSortition, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut address__ = None;
                let mut proof__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Proof => {
                            if proof__.is_some() {
                                return Err(serde::de::Error::duplicate_field("proof"));
                            }
                            proof__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(PayloadSortition {
                    address: address__.unwrap_or_default(),
                    proof: proof__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PayloadSortition", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PayloadTransfer {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.sender.is_empty() {
            len += 1;
        }
        if !self.receiver.is_empty() {
            len += 1;
        }
        if self.amount != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PayloadTransfer", len)?;
        if !self.sender.is_empty() {
            struct_ser.serialize_field("sender", &self.sender)?;
        }
        if !self.receiver.is_empty() {
            struct_ser.serialize_field("receiver", &self.receiver)?;
        }
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PayloadTransfer {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "sender",
            "receiver",
            "amount",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Sender,
            Receiver,
            Amount,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "sender" => Ok(GeneratedField::Sender),
                            "receiver" => Ok(GeneratedField::Receiver),
                            "amount" => Ok(GeneratedField::Amount),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PayloadTransfer;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PayloadTransfer")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PayloadTransfer, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut sender__ = None;
                let mut receiver__ = None;
                let mut amount__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Sender => {
                            if sender__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sender"));
                            }
                            sender__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Receiver => {
                            if receiver__.is_some() {
                                return Err(serde::de::Error::duplicate_field("receiver"));
                            }
                            receiver__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(PayloadTransfer {
                    sender: sender__.unwrap_or_default(),
                    receiver: receiver__.unwrap_or_default(),
                    amount: amount__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PayloadTransfer", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PayloadType {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "PAYLOAD_TYPE_UNSPECIFIED",
            Self::Transfer => "PAYLOAD_TYPE_TRANSFER",
            Self::Bond => "PAYLOAD_TYPE_BOND",
            Self::Sortition => "PAYLOAD_TYPE_SORTITION",
            Self::Unbond => "PAYLOAD_TYPE_UNBOND",
            Self::Withdraw => "PAYLOAD_TYPE_WITHDRAW",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for PayloadType {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "PAYLOAD_TYPE_UNSPECIFIED",
            "PAYLOAD_TYPE_TRANSFER",
            "PAYLOAD_TYPE_BOND",
            "PAYLOAD_TYPE_SORTITION",
            "PAYLOAD_TYPE_UNBOND",
            "PAYLOAD_TYPE_WITHDRAW",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PayloadType;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "PAYLOAD_TYPE_UNSPECIFIED" => Ok(PayloadType::Unspecified),
                    "PAYLOAD_TYPE_TRANSFER" => Ok(PayloadType::Transfer),
                    "PAYLOAD_TYPE_BOND" => Ok(PayloadType::Bond),
                    "PAYLOAD_TYPE_SORTITION" => Ok(PayloadType::Sortition),
                    "PAYLOAD_TYPE_UNBOND" => Ok(PayloadType::Unbond),
                    "PAYLOAD_TYPE_WITHDRAW" => Ok(PayloadType::Withdraw),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for PayloadUnbond {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.validator.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PayloadUnbond", len)?;
        if !self.validator.is_empty() {
            struct_ser.serialize_field("validator", &self.validator)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PayloadUnbond {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "validator",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Validator,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "validator" => Ok(GeneratedField::Validator),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PayloadUnbond;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PayloadUnbond")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PayloadUnbond, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut validator__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Validator => {
                            if validator__.is_some() {
                                return Err(serde::de::Error::duplicate_field("validator"));
                            }
                            validator__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(PayloadUnbond {
                    validator: validator__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PayloadUnbond", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PayloadWithdraw {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.validator_address.is_empty() {
            len += 1;
        }
        if !self.account_address.is_empty() {
            len += 1;
        }
        if self.amount != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PayloadWithdraw", len)?;
        if !self.validator_address.is_empty() {
            struct_ser.serialize_field("validatorAddress", &self.validator_address)?;
        }
        if !self.account_address.is_empty() {
            struct_ser.serialize_field("accountAddress", &self.account_address)?;
        }
        if self.amount != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("amount", ToString::to_string(&self.amount).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PayloadWithdraw {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "validator_address",
            "validatorAddress",
            "account_address",
            "accountAddress",
            "amount",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            ValidatorAddress,
            AccountAddress,
            Amount,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "validatorAddress" | "validator_address" => Ok(GeneratedField::ValidatorAddress),
                            "accountAddress" | "account_address" => Ok(GeneratedField::AccountAddress),
                            "amount" => Ok(GeneratedField::Amount),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PayloadWithdraw;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PayloadWithdraw")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PayloadWithdraw, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut validator_address__ = None;
                let mut account_address__ = None;
                let mut amount__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::ValidatorAddress => {
                            if validator_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("validatorAddress"));
                            }
                            validator_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::AccountAddress => {
                            if account_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("accountAddress"));
                            }
                            account_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Amount => {
                            if amount__.is_some() {
                                return Err(serde::de::Error::duplicate_field("amount"));
                            }
                            amount__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(PayloadWithdraw {
                    validator_address: validator_address__.unwrap_or_default(),
                    account_address: account_address__.unwrap_or_default(),
                    amount: amount__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PayloadWithdraw", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PeerInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.status != 0 {
            len += 1;
        }
        if !self.moniker.is_empty() {
            len += 1;
        }
        if !self.agent.is_empty() {
            len += 1;
        }
        if !self.peer_id.is_empty() {
            len += 1;
        }
        if !self.consensus_keys.is_empty() {
            len += 1;
        }
        if !self.consensus_addresses.is_empty() {
            len += 1;
        }
        if self.services != 0 {
            len += 1;
        }
        if !self.last_block_hash.is_empty() {
            len += 1;
        }
        if self.height != 0 {
            len += 1;
        }
        if self.last_sent != 0 {
            len += 1;
        }
        if self.last_received != 0 {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if !self.direction.is_empty() {
            len += 1;
        }
        if !self.protocols.is_empty() {
            len += 1;
        }
        if self.total_sessions != 0 {
            len += 1;
        }
        if self.completed_sessions != 0 {
            len += 1;
        }
        if self.metric_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PeerInfo", len)?;
        if self.status != 0 {
            struct_ser.serialize_field("status", &self.status)?;
        }
        if !self.moniker.is_empty() {
            struct_ser.serialize_field("moniker", &self.moniker)?;
        }
        if !self.agent.is_empty() {
            struct_ser.serialize_field("agent", &self.agent)?;
        }
        if !self.peer_id.is_empty() {
            struct_ser.serialize_field("peerId", &self.peer_id)?;
        }
        if !self.consensus_keys.is_empty() {
            struct_ser.serialize_field("consensusKeys", &self.consensus_keys)?;
        }
        if !self.consensus_addresses.is_empty() {
            struct_ser.serialize_field("consensusAddresses", &self.consensus_addresses)?;
        }
        if self.services != 0 {
            struct_ser.serialize_field("services", &self.services)?;
        }
        if !self.last_block_hash.is_empty() {
            struct_ser.serialize_field("lastBlockHash", &self.last_block_hash)?;
        }
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        if self.last_sent != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("lastSent", ToString::to_string(&self.last_sent).as_str())?;
        }
        if self.last_received != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("lastReceived", ToString::to_string(&self.last_received).as_str())?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if !self.direction.is_empty() {
            struct_ser.serialize_field("direction", &self.direction)?;
        }
        if !self.protocols.is_empty() {
            struct_ser.serialize_field("protocols", &self.protocols)?;
        }
        if self.total_sessions != 0 {
            struct_ser.serialize_field("totalSessions", &self.total_sessions)?;
        }
        if self.completed_sessions != 0 {
            struct_ser.serialize_field("completedSessions", &self.completed_sessions)?;
        }
        if let Some(v) = self.metric_info.as_ref() {
            struct_ser.serialize_field("metricInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PeerInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "status",
            "moniker",
            "agent",
            "peer_id",
            "peerId",
            "consensus_keys",
            "consensusKeys",
            "consensus_addresses",
            "consensusAddresses",
            "services",
            "last_block_hash",
            "lastBlockHash",
            "height",
            "last_sent",
            "lastSent",
            "last_received",
            "lastReceived",
            "address",
            "direction",
            "protocols",
            "total_sessions",
            "totalSessions",
            "completed_sessions",
            "completedSessions",
            "metric_info",
            "metricInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Status,
            Moniker,
            Agent,
            PeerId,
            ConsensusKeys,
            ConsensusAddresses,
            Services,
            LastBlockHash,
            Height,
            LastSent,
            LastReceived,
            Address,
            Direction,
            Protocols,
            TotalSessions,
            CompletedSessions,
            MetricInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "status" => Ok(GeneratedField::Status),
                            "moniker" => Ok(GeneratedField::Moniker),
                            "agent" => Ok(GeneratedField::Agent),
                            "peerId" | "peer_id" => Ok(GeneratedField::PeerId),
                            "consensusKeys" | "consensus_keys" => Ok(GeneratedField::ConsensusKeys),
                            "consensusAddresses" | "consensus_addresses" => Ok(GeneratedField::ConsensusAddresses),
                            "services" => Ok(GeneratedField::Services),
                            "lastBlockHash" | "last_block_hash" => Ok(GeneratedField::LastBlockHash),
                            "height" => Ok(GeneratedField::Height),
                            "lastSent" | "last_sent" => Ok(GeneratedField::LastSent),
                            "lastReceived" | "last_received" => Ok(GeneratedField::LastReceived),
                            "address" => Ok(GeneratedField::Address),
                            "direction" => Ok(GeneratedField::Direction),
                            "protocols" => Ok(GeneratedField::Protocols),
                            "totalSessions" | "total_sessions" => Ok(GeneratedField::TotalSessions),
                            "completedSessions" | "completed_sessions" => Ok(GeneratedField::CompletedSessions),
                            "metricInfo" | "metric_info" => Ok(GeneratedField::MetricInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PeerInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PeerInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PeerInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut status__ = None;
                let mut moniker__ = None;
                let mut agent__ = None;
                let mut peer_id__ = None;
                let mut consensus_keys__ = None;
                let mut consensus_addresses__ = None;
                let mut services__ = None;
                let mut last_block_hash__ = None;
                let mut height__ = None;
                let mut last_sent__ = None;
                let mut last_received__ = None;
                let mut address__ = None;
                let mut direction__ = None;
                let mut protocols__ = None;
                let mut total_sessions__ = None;
                let mut completed_sessions__ = None;
                let mut metric_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Status => {
                            if status__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            status__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Moniker => {
                            if moniker__.is_some() {
                                return Err(serde::de::Error::duplicate_field("moniker"));
                            }
                            moniker__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Agent => {
                            if agent__.is_some() {
                                return Err(serde::de::Error::duplicate_field("agent"));
                            }
                            agent__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PeerId => {
                            if peer_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("peerId"));
                            }
                            peer_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ConsensusKeys => {
                            if consensus_keys__.is_some() {
                                return Err(serde::de::Error::duplicate_field("consensusKeys"));
                            }
                            consensus_keys__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ConsensusAddresses => {
                            if consensus_addresses__.is_some() {
                                return Err(serde::de::Error::duplicate_field("consensusAddresses"));
                            }
                            consensus_addresses__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Services => {
                            if services__.is_some() {
                                return Err(serde::de::Error::duplicate_field("services"));
                            }
                            services__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastBlockHash => {
                            if last_block_hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastBlockHash"));
                            }
                            last_block_hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastSent => {
                            if last_sent__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastSent"));
                            }
                            last_sent__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastReceived => {
                            if last_received__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastReceived"));
                            }
                            last_received__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Direction => {
                            if direction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("direction"));
                            }
                            direction__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Protocols => {
                            if protocols__.is_some() {
                                return Err(serde::de::Error::duplicate_field("protocols"));
                            }
                            protocols__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TotalSessions => {
                            if total_sessions__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalSessions"));
                            }
                            total_sessions__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::CompletedSessions => {
                            if completed_sessions__.is_some() {
                                return Err(serde::de::Error::duplicate_field("completedSessions"));
                            }
                            completed_sessions__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MetricInfo => {
                            if metric_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metricInfo"));
                            }
                            metric_info__ = map_.next_value()?;
                        }
                    }
                }
                Ok(PeerInfo {
                    status: status__.unwrap_or_default(),
                    moniker: moniker__.unwrap_or_default(),
                    agent: agent__.unwrap_or_default(),
                    peer_id: peer_id__.unwrap_or_default(),
                    consensus_keys: consensus_keys__.unwrap_or_default(),
                    consensus_addresses: consensus_addresses__.unwrap_or_default(),
                    services: services__.unwrap_or_default(),
                    last_block_hash: last_block_hash__.unwrap_or_default(),
                    height: height__.unwrap_or_default(),
                    last_sent: last_sent__.unwrap_or_default(),
                    last_received: last_received__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    direction: direction__.unwrap_or_default(),
                    protocols: protocols__.unwrap_or_default(),
                    total_sessions: total_sessions__.unwrap_or_default(),
                    completed_sessions: completed_sessions__.unwrap_or_default(),
                    metric_info: metric_info__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.PeerInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ProposalInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.height != 0 {
            len += 1;
        }
        if self.round != 0 {
            len += 1;
        }
        if !self.block_data.is_empty() {
            len += 1;
        }
        if !self.signature.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ProposalInfo", len)?;
        if self.height != 0 {
            struct_ser.serialize_field("height", &self.height)?;
        }
        if self.round != 0 {
            struct_ser.serialize_field("round", &self.round)?;
        }
        if !self.block_data.is_empty() {
            struct_ser.serialize_field("blockData", &self.block_data)?;
        }
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ProposalInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "height",
            "round",
            "block_data",
            "blockData",
            "signature",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Height,
            Round,
            BlockData,
            Signature,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "height" => Ok(GeneratedField::Height),
                            "round" => Ok(GeneratedField::Round),
                            "blockData" | "block_data" => Ok(GeneratedField::BlockData),
                            "signature" => Ok(GeneratedField::Signature),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ProposalInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ProposalInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ProposalInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut height__ = None;
                let mut round__ = None;
                let mut block_data__ = None;
                let mut signature__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Height => {
                            if height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("height"));
                            }
                            height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Round => {
                            if round__.is_some() {
                                return Err(serde::de::Error::duplicate_field("round"));
                            }
                            round__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::BlockData => {
                            if block_data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("blockData"));
                            }
                            block_data__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(ProposalInfo {
                    height: height__.unwrap_or_default(),
                    round: round__.unwrap_or_default(),
                    block_data: block_data__.unwrap_or_default(),
                    signature: signature__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ProposalInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PublicKeyAggregationRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.public_keys.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PublicKeyAggregationRequest", len)?;
        if !self.public_keys.is_empty() {
            struct_ser.serialize_field("publicKeys", &self.public_keys)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PublicKeyAggregationRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "public_keys",
            "publicKeys",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PublicKeys,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "publicKeys" | "public_keys" => Ok(GeneratedField::PublicKeys),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PublicKeyAggregationRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PublicKeyAggregationRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PublicKeyAggregationRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut public_keys__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PublicKeys => {
                            if public_keys__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKeys"));
                            }
                            public_keys__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(PublicKeyAggregationRequest {
                    public_keys: public_keys__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PublicKeyAggregationRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for PublicKeyAggregationResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.public_key.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.PublicKeyAggregationResponse", len)?;
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PublicKeyAggregationResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "public_key",
            "publicKey",
            "address",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PublicKey,
            Address,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            "address" => Ok(GeneratedField::Address),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PublicKeyAggregationResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.PublicKeyAggregationResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PublicKeyAggregationResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut public_key__ = None;
                let mut address__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(PublicKeyAggregationResponse {
                    public_key: public_key__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.PublicKeyAggregationResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RestoreWalletRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.mnemonic.is_empty() {
            len += 1;
        }
        if !self.password.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.RestoreWalletRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.mnemonic.is_empty() {
            struct_ser.serialize_field("mnemonic", &self.mnemonic)?;
        }
        if !self.password.is_empty() {
            struct_ser.serialize_field("password", &self.password)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RestoreWalletRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "mnemonic",
            "password",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Mnemonic,
            Password,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "mnemonic" => Ok(GeneratedField::Mnemonic),
                            "password" => Ok(GeneratedField::Password),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RestoreWalletRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.RestoreWalletRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RestoreWalletRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut mnemonic__ = None;
                let mut password__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Mnemonic => {
                            if mnemonic__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mnemonic"));
                            }
                            mnemonic__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Password => {
                            if password__.is_some() {
                                return Err(serde::de::Error::duplicate_field("password"));
                            }
                            password__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RestoreWalletRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    mnemonic: mnemonic__.unwrap_or_default(),
                    password: password__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.RestoreWalletRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RestoreWalletResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.RestoreWalletResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RestoreWalletResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RestoreWalletResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.RestoreWalletResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RestoreWalletResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RestoreWalletResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.RestoreWalletResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SetAddressLabelRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.password.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if !self.label.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SetAddressLabelRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.password.is_empty() {
            struct_ser.serialize_field("password", &self.password)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if !self.label.is_empty() {
            struct_ser.serialize_field("label", &self.label)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SetAddressLabelRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "password",
            "address",
            "label",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Password,
            Address,
            Label,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "password" => Ok(GeneratedField::Password),
                            "address" => Ok(GeneratedField::Address),
                            "label" => Ok(GeneratedField::Label),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SetAddressLabelRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SetAddressLabelRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SetAddressLabelRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut password__ = None;
                let mut address__ = None;
                let mut label__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Password => {
                            if password__.is_some() {
                                return Err(serde::de::Error::duplicate_field("password"));
                            }
                            password__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Label => {
                            if label__.is_some() {
                                return Err(serde::de::Error::duplicate_field("label"));
                            }
                            label__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SetAddressLabelRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    password: password__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    label: label__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SetAddressLabelRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SetAddressLabelResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("pactus.SetAddressLabelResponse", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SetAddressLabelResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SetAddressLabelResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SetAddressLabelResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SetAddressLabelResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(SetAddressLabelResponse {
                })
            }
        }
        deserializer.deserialize_struct("pactus.SetAddressLabelResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignMessageRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.password.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if !self.message.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignMessageRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.password.is_empty() {
            struct_ser.serialize_field("password", &self.password)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignMessageRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "password",
            "address",
            "message",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            Password,
            Address,
            Message,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "password" => Ok(GeneratedField::Password),
                            "address" => Ok(GeneratedField::Address),
                            "message" => Ok(GeneratedField::Message),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignMessageRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignMessageRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignMessageRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut password__ = None;
                let mut address__ = None;
                let mut message__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Password => {
                            if password__.is_some() {
                                return Err(serde::de::Error::duplicate_field("password"));
                            }
                            password__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignMessageRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    password: password__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    message: message__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignMessageRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignMessageResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.signature.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignMessageResponse", len)?;
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignMessageResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "signature",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Signature,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "signature" => Ok(GeneratedField::Signature),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignMessageResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignMessageResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignMessageResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut signature__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignMessageResponse {
                    signature: signature__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignMessageResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignMessageWithPrivateKeyRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.private_key.is_empty() {
            len += 1;
        }
        if !self.message.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignMessageWithPrivateKeyRequest", len)?;
        if !self.private_key.is_empty() {
            struct_ser.serialize_field("privateKey", &self.private_key)?;
        }
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignMessageWithPrivateKeyRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "private_key",
            "privateKey",
            "message",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PrivateKey,
            Message,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "privateKey" | "private_key" => Ok(GeneratedField::PrivateKey),
                            "message" => Ok(GeneratedField::Message),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignMessageWithPrivateKeyRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignMessageWithPrivateKeyRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignMessageWithPrivateKeyRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut private_key__ = None;
                let mut message__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PrivateKey => {
                            if private_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("privateKey"));
                            }
                            private_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignMessageWithPrivateKeyRequest {
                    private_key: private_key__.unwrap_or_default(),
                    message: message__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignMessageWithPrivateKeyRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignMessageWithPrivateKeyResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.signature.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignMessageWithPrivateKeyResponse", len)?;
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignMessageWithPrivateKeyResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "signature",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Signature,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "signature" => Ok(GeneratedField::Signature),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignMessageWithPrivateKeyResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignMessageWithPrivateKeyResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignMessageWithPrivateKeyResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut signature__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignMessageWithPrivateKeyResponse {
                    signature: signature__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignMessageWithPrivateKeyResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignRawTransactionRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        if !self.raw_transaction.is_empty() {
            len += 1;
        }
        if !self.password.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignRawTransactionRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        if !self.raw_transaction.is_empty() {
            struct_ser.serialize_field("rawTransaction", &self.raw_transaction)?;
        }
        if !self.password.is_empty() {
            struct_ser.serialize_field("password", &self.password)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignRawTransactionRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
            "raw_transaction",
            "rawTransaction",
            "password",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
            RawTransaction,
            Password,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            "rawTransaction" | "raw_transaction" => Ok(GeneratedField::RawTransaction),
                            "password" => Ok(GeneratedField::Password),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignRawTransactionRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignRawTransactionRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignRawTransactionRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                let mut raw_transaction__ = None;
                let mut password__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::RawTransaction => {
                            if raw_transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("rawTransaction"));
                            }
                            raw_transaction__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Password => {
                            if password__.is_some() {
                                return Err(serde::de::Error::duplicate_field("password"));
                            }
                            password__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignRawTransactionRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                    raw_transaction: raw_transaction__.unwrap_or_default(),
                    password: password__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignRawTransactionRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignRawTransactionResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.transaction_id.is_empty() {
            len += 1;
        }
        if !self.signed_raw_transaction.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignRawTransactionResponse", len)?;
        if !self.transaction_id.is_empty() {
            struct_ser.serialize_field("transactionId", &self.transaction_id)?;
        }
        if !self.signed_raw_transaction.is_empty() {
            struct_ser.serialize_field("signedRawTransaction", &self.signed_raw_transaction)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignRawTransactionResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "transaction_id",
            "transactionId",
            "signed_raw_transaction",
            "signedRawTransaction",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TransactionId,
            SignedRawTransaction,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "transactionId" | "transaction_id" => Ok(GeneratedField::TransactionId),
                            "signedRawTransaction" | "signed_raw_transaction" => Ok(GeneratedField::SignedRawTransaction),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignRawTransactionResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignRawTransactionResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignRawTransactionResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut transaction_id__ = None;
                let mut signed_raw_transaction__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TransactionId => {
                            if transaction_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("transactionId"));
                            }
                            transaction_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SignedRawTransaction => {
                            if signed_raw_transaction__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signedRawTransaction"));
                            }
                            signed_raw_transaction__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignRawTransactionResponse {
                    transaction_id: transaction_id__.unwrap_or_default(),
                    signed_raw_transaction: signed_raw_transaction__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignRawTransactionResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignatureAggregationRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.signatures.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignatureAggregationRequest", len)?;
        if !self.signatures.is_empty() {
            struct_ser.serialize_field("signatures", &self.signatures)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignatureAggregationRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "signatures",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Signatures,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "signatures" => Ok(GeneratedField::Signatures),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignatureAggregationRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignatureAggregationRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignatureAggregationRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut signatures__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Signatures => {
                            if signatures__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signatures"));
                            }
                            signatures__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignatureAggregationRequest {
                    signatures: signatures__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignatureAggregationRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SignatureAggregationResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.signature.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.SignatureAggregationResponse", len)?;
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SignatureAggregationResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "signature",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Signature,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "signature" => Ok(GeneratedField::Signature),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SignatureAggregationResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.SignatureAggregationResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SignatureAggregationResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut signature__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SignatureAggregationResponse {
                    signature: signature__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.SignatureAggregationResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TransactionInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if !self.data.is_empty() {
            len += 1;
        }
        if self.version != 0 {
            len += 1;
        }
        if self.lock_time != 0 {
            len += 1;
        }
        if self.value != 0 {
            len += 1;
        }
        if self.fee != 0 {
            len += 1;
        }
        if self.payload_type != 0 {
            len += 1;
        }
        if !self.memo.is_empty() {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        if !self.signature.is_empty() {
            len += 1;
        }
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.TransactionInfo", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.data.is_empty() {
            struct_ser.serialize_field("data", &self.data)?;
        }
        if self.version != 0 {
            struct_ser.serialize_field("version", &self.version)?;
        }
        if self.lock_time != 0 {
            struct_ser.serialize_field("lockTime", &self.lock_time)?;
        }
        if self.value != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("value", ToString::to_string(&self.value).as_str())?;
        }
        if self.fee != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("fee", ToString::to_string(&self.fee).as_str())?;
        }
        if self.payload_type != 0 {
            let v = PayloadType::try_from(self.payload_type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.payload_type)))?;
            struct_ser.serialize_field("payloadType", &v)?;
        }
        if !self.memo.is_empty() {
            struct_ser.serialize_field("memo", &self.memo)?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        if let Some(v) = self.payload.as_ref() {
            match v {
                transaction_info::Payload::Transfer(v) => {
                    struct_ser.serialize_field("transfer", v)?;
                }
                transaction_info::Payload::Bond(v) => {
                    struct_ser.serialize_field("bond", v)?;
                }
                transaction_info::Payload::Sortition(v) => {
                    struct_ser.serialize_field("sortition", v)?;
                }
                transaction_info::Payload::Unbond(v) => {
                    struct_ser.serialize_field("unbond", v)?;
                }
                transaction_info::Payload::Withdraw(v) => {
                    struct_ser.serialize_field("withdraw", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TransactionInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "data",
            "version",
            "lock_time",
            "lockTime",
            "value",
            "fee",
            "payload_type",
            "payloadType",
            "memo",
            "public_key",
            "publicKey",
            "signature",
            "transfer",
            "bond",
            "sortition",
            "unbond",
            "withdraw",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Data,
            Version,
            LockTime,
            Value,
            Fee,
            PayloadType,
            Memo,
            PublicKey,
            Signature,
            Transfer,
            Bond,
            Sortition,
            Unbond,
            Withdraw,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "id" => Ok(GeneratedField::Id),
                            "data" => Ok(GeneratedField::Data),
                            "version" => Ok(GeneratedField::Version),
                            "lockTime" | "lock_time" => Ok(GeneratedField::LockTime),
                            "value" => Ok(GeneratedField::Value),
                            "fee" => Ok(GeneratedField::Fee),
                            "payloadType" | "payload_type" => Ok(GeneratedField::PayloadType),
                            "memo" => Ok(GeneratedField::Memo),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            "signature" => Ok(GeneratedField::Signature),
                            "transfer" => Ok(GeneratedField::Transfer),
                            "bond" => Ok(GeneratedField::Bond),
                            "sortition" => Ok(GeneratedField::Sortition),
                            "unbond" => Ok(GeneratedField::Unbond),
                            "withdraw" => Ok(GeneratedField::Withdraw),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TransactionInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.TransactionInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TransactionInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut data__ = None;
                let mut version__ = None;
                let mut lock_time__ = None;
                let mut value__ = None;
                let mut fee__ = None;
                let mut payload_type__ = None;
                let mut memo__ = None;
                let mut public_key__ = None;
                let mut signature__ = None;
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Data => {
                            if data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("data"));
                            }
                            data__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LockTime => {
                            if lock_time__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockTime"));
                            }
                            lock_time__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Fee => {
                            if fee__.is_some() {
                                return Err(serde::de::Error::duplicate_field("fee"));
                            }
                            fee__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PayloadType => {
                            if payload_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("payloadType"));
                            }
                            payload_type__ = Some(map_.next_value::<PayloadType>()? as i32);
                        }
                        GeneratedField::Memo => {
                            if memo__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memo"));
                            }
                            memo__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Transfer => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("transfer"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(transaction_info::Payload::Transfer)
;
                        }
                        GeneratedField::Bond => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("bond"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(transaction_info::Payload::Bond)
;
                        }
                        GeneratedField::Sortition => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sortition"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(transaction_info::Payload::Sortition)
;
                        }
                        GeneratedField::Unbond => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("unbond"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(transaction_info::Payload::Unbond)
;
                        }
                        GeneratedField::Withdraw => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("withdraw"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(transaction_info::Payload::Withdraw)
;
                        }
                    }
                }
                Ok(TransactionInfo {
                    id: id__.unwrap_or_default(),
                    data: data__.unwrap_or_default(),
                    version: version__.unwrap_or_default(),
                    lock_time: lock_time__.unwrap_or_default(),
                    value: value__.unwrap_or_default(),
                    fee: fee__.unwrap_or_default(),
                    payload_type: payload_type__.unwrap_or_default(),
                    memo: memo__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                    signature: signature__.unwrap_or_default(),
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("pactus.TransactionInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TransactionVerbosity {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Data => "TRANSACTION_VERBOSITY_DATA",
            Self::Info => "TRANSACTION_VERBOSITY_INFO",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for TransactionVerbosity {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "TRANSACTION_VERBOSITY_DATA",
            "TRANSACTION_VERBOSITY_INFO",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TransactionVerbosity;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "TRANSACTION_VERBOSITY_DATA" => Ok(TransactionVerbosity::Data),
                    "TRANSACTION_VERBOSITY_INFO" => Ok(TransactionVerbosity::Info),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for UnloadWalletRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.UnloadWalletRequest", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UnloadWalletRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UnloadWalletRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.UnloadWalletRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UnloadWalletRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(UnloadWalletRequest {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.UnloadWalletRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UnloadWalletResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.wallet_name.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.UnloadWalletResponse", len)?;
        if !self.wallet_name.is_empty() {
            struct_ser.serialize_field("walletName", &self.wallet_name)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UnloadWalletResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wallet_name",
            "walletName",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WalletName,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "walletName" | "wallet_name" => Ok(GeneratedField::WalletName),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = UnloadWalletResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.UnloadWalletResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UnloadWalletResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut wallet_name__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WalletName => {
                            if wallet_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("walletName"));
                            }
                            wallet_name__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(UnloadWalletResponse {
                    wallet_name: wallet_name__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.UnloadWalletResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ValidatorInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.hash.is_empty() {
            len += 1;
        }
        if !self.data.is_empty() {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        if self.number != 0 {
            len += 1;
        }
        if self.stake != 0 {
            len += 1;
        }
        if self.last_bonding_height != 0 {
            len += 1;
        }
        if self.last_sortition_height != 0 {
            len += 1;
        }
        if self.unbonding_height != 0 {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if self.availability_score != 0. {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ValidatorInfo", len)?;
        if !self.hash.is_empty() {
            struct_ser.serialize_field("hash", &self.hash)?;
        }
        if !self.data.is_empty() {
            struct_ser.serialize_field("data", &self.data)?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        if self.number != 0 {
            struct_ser.serialize_field("number", &self.number)?;
        }
        if self.stake != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("stake", ToString::to_string(&self.stake).as_str())?;
        }
        if self.last_bonding_height != 0 {
            struct_ser.serialize_field("lastBondingHeight", &self.last_bonding_height)?;
        }
        if self.last_sortition_height != 0 {
            struct_ser.serialize_field("lastSortitionHeight", &self.last_sortition_height)?;
        }
        if self.unbonding_height != 0 {
            struct_ser.serialize_field("unbondingHeight", &self.unbonding_height)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if self.availability_score != 0. {
            struct_ser.serialize_field("availabilityScore", &self.availability_score)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ValidatorInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "hash",
            "data",
            "public_key",
            "publicKey",
            "number",
            "stake",
            "last_bonding_height",
            "lastBondingHeight",
            "last_sortition_height",
            "lastSortitionHeight",
            "unbonding_height",
            "unbondingHeight",
            "address",
            "availability_score",
            "availabilityScore",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Hash,
            Data,
            PublicKey,
            Number,
            Stake,
            LastBondingHeight,
            LastSortitionHeight,
            UnbondingHeight,
            Address,
            AvailabilityScore,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "hash" => Ok(GeneratedField::Hash),
                            "data" => Ok(GeneratedField::Data),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            "number" => Ok(GeneratedField::Number),
                            "stake" => Ok(GeneratedField::Stake),
                            "lastBondingHeight" | "last_bonding_height" => Ok(GeneratedField::LastBondingHeight),
                            "lastSortitionHeight" | "last_sortition_height" => Ok(GeneratedField::LastSortitionHeight),
                            "unbondingHeight" | "unbonding_height" => Ok(GeneratedField::UnbondingHeight),
                            "address" => Ok(GeneratedField::Address),
                            "availabilityScore" | "availability_score" => Ok(GeneratedField::AvailabilityScore),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ValidatorInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ValidatorInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ValidatorInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut hash__ = None;
                let mut data__ = None;
                let mut public_key__ = None;
                let mut number__ = None;
                let mut stake__ = None;
                let mut last_bonding_height__ = None;
                let mut last_sortition_height__ = None;
                let mut unbonding_height__ = None;
                let mut address__ = None;
                let mut availability_score__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Hash => {
                            if hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hash"));
                            }
                            hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Data => {
                            if data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("data"));
                            }
                            data__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Number => {
                            if number__.is_some() {
                                return Err(serde::de::Error::duplicate_field("number"));
                            }
                            number__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Stake => {
                            if stake__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stake"));
                            }
                            stake__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastBondingHeight => {
                            if last_bonding_height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastBondingHeight"));
                            }
                            last_bonding_height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LastSortitionHeight => {
                            if last_sortition_height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastSortitionHeight"));
                            }
                            last_sortition_height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::UnbondingHeight => {
                            if unbonding_height__.is_some() {
                                return Err(serde::de::Error::duplicate_field("unbondingHeight"));
                            }
                            unbonding_height__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::AvailabilityScore => {
                            if availability_score__.is_some() {
                                return Err(serde::de::Error::duplicate_field("availabilityScore"));
                            }
                            availability_score__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(ValidatorInfo {
                    hash: hash__.unwrap_or_default(),
                    data: data__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                    number: number__.unwrap_or_default(),
                    stake: stake__.unwrap_or_default(),
                    last_bonding_height: last_bonding_height__.unwrap_or_default(),
                    last_sortition_height: last_sortition_height__.unwrap_or_default(),
                    unbonding_height: unbonding_height__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    availability_score: availability_score__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ValidatorInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for VerifyMessageRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.message.is_empty() {
            len += 1;
        }
        if !self.signature.is_empty() {
            len += 1;
        }
        if !self.public_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.VerifyMessageRequest", len)?;
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        if !self.signature.is_empty() {
            struct_ser.serialize_field("signature", &self.signature)?;
        }
        if !self.public_key.is_empty() {
            struct_ser.serialize_field("publicKey", &self.public_key)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for VerifyMessageRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "message",
            "signature",
            "public_key",
            "publicKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Message,
            Signature,
            PublicKey,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "message" => Ok(GeneratedField::Message),
                            "signature" => Ok(GeneratedField::Signature),
                            "publicKey" | "public_key" => Ok(GeneratedField::PublicKey),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = VerifyMessageRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.VerifyMessageRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<VerifyMessageRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut message__ = None;
                let mut signature__ = None;
                let mut public_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Signature => {
                            if signature__.is_some() {
                                return Err(serde::de::Error::duplicate_field("signature"));
                            }
                            signature__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PublicKey => {
                            if public_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("publicKey"));
                            }
                            public_key__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(VerifyMessageRequest {
                    message: message__.unwrap_or_default(),
                    signature: signature__.unwrap_or_default(),
                    public_key: public_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.VerifyMessageRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for VerifyMessageResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.is_valid {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.VerifyMessageResponse", len)?;
        if self.is_valid {
            struct_ser.serialize_field("isValid", &self.is_valid)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for VerifyMessageResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "is_valid",
            "isValid",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            IsValid,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "isValid" | "is_valid" => Ok(GeneratedField::IsValid),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = VerifyMessageResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.VerifyMessageResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<VerifyMessageResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut is_valid__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::IsValid => {
                            if is_valid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isValid"));
                            }
                            is_valid__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(VerifyMessageResponse {
                    is_valid: is_valid__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.VerifyMessageResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for VoteInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.r#type != 0 {
            len += 1;
        }
        if !self.voter.is_empty() {
            len += 1;
        }
        if !self.block_hash.is_empty() {
            len += 1;
        }
        if self.round != 0 {
            len += 1;
        }
        if self.cp_round != 0 {
            len += 1;
        }
        if self.cp_value != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.VoteInfo", len)?;
        if self.r#type != 0 {
            let v = VoteType::try_from(self.r#type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.r#type)))?;
            struct_ser.serialize_field("type", &v)?;
        }
        if !self.voter.is_empty() {
            struct_ser.serialize_field("voter", &self.voter)?;
        }
        if !self.block_hash.is_empty() {
            struct_ser.serialize_field("blockHash", &self.block_hash)?;
        }
        if self.round != 0 {
            struct_ser.serialize_field("round", &self.round)?;
        }
        if self.cp_round != 0 {
            struct_ser.serialize_field("cpRound", &self.cp_round)?;
        }
        if self.cp_value != 0 {
            struct_ser.serialize_field("cpValue", &self.cp_value)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for VoteInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "type",
            "voter",
            "block_hash",
            "blockHash",
            "round",
            "cp_round",
            "cpRound",
            "cp_value",
            "cpValue",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Type,
            Voter,
            BlockHash,
            Round,
            CpRound,
            CpValue,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "type" => Ok(GeneratedField::Type),
                            "voter" => Ok(GeneratedField::Voter),
                            "blockHash" | "block_hash" => Ok(GeneratedField::BlockHash),
                            "round" => Ok(GeneratedField::Round),
                            "cpRound" | "cp_round" => Ok(GeneratedField::CpRound),
                            "cpValue" | "cp_value" => Ok(GeneratedField::CpValue),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = VoteInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.VoteInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<VoteInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#type__ = None;
                let mut voter__ = None;
                let mut block_hash__ = None;
                let mut round__ = None;
                let mut cp_round__ = None;
                let mut cp_value__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value::<VoteType>()? as i32);
                        }
                        GeneratedField::Voter => {
                            if voter__.is_some() {
                                return Err(serde::de::Error::duplicate_field("voter"));
                            }
                            voter__ = Some(map_.next_value()?);
                        }
                        GeneratedField::BlockHash => {
                            if block_hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("blockHash"));
                            }
                            block_hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Round => {
                            if round__.is_some() {
                                return Err(serde::de::Error::duplicate_field("round"));
                            }
                            round__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::CpRound => {
                            if cp_round__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cpRound"));
                            }
                            cp_round__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::CpValue => {
                            if cp_value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cpValue"));
                            }
                            cp_value__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(VoteInfo {
                    r#type: r#type__.unwrap_or_default(),
                    voter: voter__.unwrap_or_default(),
                    block_hash: block_hash__.unwrap_or_default(),
                    round: round__.unwrap_or_default(),
                    cp_round: cp_round__.unwrap_or_default(),
                    cp_value: cp_value__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.VoteInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for VoteType {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unspecified => "VOTE_TYPE_UNSPECIFIED",
            Self::Prepare => "VOTE_TYPE_PREPARE",
            Self::Precommit => "VOTE_TYPE_PRECOMMIT",
            Self::CpPreVote => "VOTE_TYPE_CP_PRE_VOTE",
            Self::CpMainVote => "VOTE_TYPE_CP_MAIN_VOTE",
            Self::CpDecided => "VOTE_TYPE_CP_DECIDED",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for VoteType {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "VOTE_TYPE_UNSPECIFIED",
            "VOTE_TYPE_PREPARE",
            "VOTE_TYPE_PRECOMMIT",
            "VOTE_TYPE_CP_PRE_VOTE",
            "VOTE_TYPE_CP_MAIN_VOTE",
            "VOTE_TYPE_CP_DECIDED",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = VoteType;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "VOTE_TYPE_UNSPECIFIED" => Ok(VoteType::Unspecified),
                    "VOTE_TYPE_PREPARE" => Ok(VoteType::Prepare),
                    "VOTE_TYPE_PRECOMMIT" => Ok(VoteType::Precommit),
                    "VOTE_TYPE_CP_PRE_VOTE" => Ok(VoteType::CpPreVote),
                    "VOTE_TYPE_CP_MAIN_VOTE" => Ok(VoteType::CpMainVote),
                    "VOTE_TYPE_CP_DECIDED" => Ok(VoteType::CpDecided),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for ZmqPublisherInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.topic.is_empty() {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if self.hwm != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pactus.ZMQPublisherInfo", len)?;
        if !self.topic.is_empty() {
            struct_ser.serialize_field("topic", &self.topic)?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if self.hwm != 0 {
            struct_ser.serialize_field("hwm", &self.hwm)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ZmqPublisherInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "topic",
            "address",
            "hwm",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Topic,
            Address,
            Hwm,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "topic" => Ok(GeneratedField::Topic),
                            "address" => Ok(GeneratedField::Address),
                            "hwm" => Ok(GeneratedField::Hwm),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ZmqPublisherInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pactus.ZMQPublisherInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ZmqPublisherInfo, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut topic__ = None;
                let mut address__ = None;
                let mut hwm__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Topic => {
                            if topic__.is_some() {
                                return Err(serde::de::Error::duplicate_field("topic"));
                            }
                            topic__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Hwm => {
                            if hwm__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hwm"));
                            }
                            hwm__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(ZmqPublisherInfo {
                    topic: topic__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    hwm: hwm__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pactus.ZMQPublisherInfo", FIELDS, GeneratedVisitor)
    }
}
