// This file is generated. Do not edit
// @generated

// https://github.com/Manishearth/rust-clippy/issues/702
#![allow(unknown_lints)]
#![allow(clippy::all)]

#![allow(box_pointers)]
#![allow(dead_code)]
#![allow(missing_docs)]
#![allow(non_camel_case_types)]
#![allow(non_snake_case)]
#![allow(non_upper_case_globals)]
#![allow(trivial_casts)]
#![allow(unsafe_code)]
#![allow(unused_imports)]
#![allow(unused_results)]

const METHOD_WALLET_CREATE_WALLET: ::grpcio::Method<super::wallet::CreateWalletRequest, super::wallet::CreateWalletResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/CreateWallet",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_WALLET_LOAD_WALLET: ::grpcio::Method<super::wallet::LoadWalletRequest, super::wallet::LoadWalletResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/LoadWallet",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_WALLET_UNLOAD_WALLET: ::grpcio::Method<super::wallet::UnloadWalletRequest, super::wallet::UnloadWalletResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/UnloadWallet",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_WALLET_LOCK_WALLET: ::grpcio::Method<super::wallet::LockWalletRequest, super::wallet::LockWalletResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/LockWallet",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_WALLET_UNLOCK_WALLET: ::grpcio::Method<super::wallet::UnlockWalletRequest, super::wallet::UnlockWalletResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/UnlockWallet",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_WALLET_SIGN_RAW_TRANSACTION: ::grpcio::Method<super::wallet::SignRawTransactionRequest, super::wallet::SignRawTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/SignRawTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_WALLET_GET_VALIDATOR_ADDRESS: ::grpcio::Method<super::wallet::GetValidatorAddressRequest, super::wallet::GetValidatorAddressResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Wallet/GetValidatorAddress",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

#[derive(Clone)]
pub struct WalletClient {
    pub client: ::grpcio::Client,
}

impl WalletClient {
    pub fn new(channel: ::grpcio::Channel) -> Self {
        WalletClient {
            client: ::grpcio::Client::new(channel),
        }
    }

    pub fn create_wallet_opt(&self, req: &super::wallet::CreateWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::CreateWalletResponse> {
        self.client.unary_call(&METHOD_WALLET_CREATE_WALLET, req, opt)
    }

    pub fn create_wallet(&self, req: &super::wallet::CreateWalletRequest) -> ::grpcio::Result<super::wallet::CreateWalletResponse> {
        self.create_wallet_opt(req, ::grpcio::CallOption::default())
    }

    pub fn create_wallet_async_opt(&self, req: &super::wallet::CreateWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::CreateWalletResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_CREATE_WALLET, req, opt)
    }

    pub fn create_wallet_async(&self, req: &super::wallet::CreateWalletRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::CreateWalletResponse>> {
        self.create_wallet_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn load_wallet_opt(&self, req: &super::wallet::LoadWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::LoadWalletResponse> {
        self.client.unary_call(&METHOD_WALLET_LOAD_WALLET, req, opt)
    }

    pub fn load_wallet(&self, req: &super::wallet::LoadWalletRequest) -> ::grpcio::Result<super::wallet::LoadWalletResponse> {
        self.load_wallet_opt(req, ::grpcio::CallOption::default())
    }

    pub fn load_wallet_async_opt(&self, req: &super::wallet::LoadWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::LoadWalletResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_LOAD_WALLET, req, opt)
    }

    pub fn load_wallet_async(&self, req: &super::wallet::LoadWalletRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::LoadWalletResponse>> {
        self.load_wallet_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn unload_wallet_opt(&self, req: &super::wallet::UnloadWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::UnloadWalletResponse> {
        self.client.unary_call(&METHOD_WALLET_UNLOAD_WALLET, req, opt)
    }

    pub fn unload_wallet(&self, req: &super::wallet::UnloadWalletRequest) -> ::grpcio::Result<super::wallet::UnloadWalletResponse> {
        self.unload_wallet_opt(req, ::grpcio::CallOption::default())
    }

    pub fn unload_wallet_async_opt(&self, req: &super::wallet::UnloadWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::UnloadWalletResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_UNLOAD_WALLET, req, opt)
    }

    pub fn unload_wallet_async(&self, req: &super::wallet::UnloadWalletRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::UnloadWalletResponse>> {
        self.unload_wallet_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn lock_wallet_opt(&self, req: &super::wallet::LockWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::LockWalletResponse> {
        self.client.unary_call(&METHOD_WALLET_LOCK_WALLET, req, opt)
    }

    pub fn lock_wallet(&self, req: &super::wallet::LockWalletRequest) -> ::grpcio::Result<super::wallet::LockWalletResponse> {
        self.lock_wallet_opt(req, ::grpcio::CallOption::default())
    }

    pub fn lock_wallet_async_opt(&self, req: &super::wallet::LockWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::LockWalletResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_LOCK_WALLET, req, opt)
    }

    pub fn lock_wallet_async(&self, req: &super::wallet::LockWalletRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::LockWalletResponse>> {
        self.lock_wallet_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn unlock_wallet_opt(&self, req: &super::wallet::UnlockWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::UnlockWalletResponse> {
        self.client.unary_call(&METHOD_WALLET_UNLOCK_WALLET, req, opt)
    }

    pub fn unlock_wallet(&self, req: &super::wallet::UnlockWalletRequest) -> ::grpcio::Result<super::wallet::UnlockWalletResponse> {
        self.unlock_wallet_opt(req, ::grpcio::CallOption::default())
    }

    pub fn unlock_wallet_async_opt(&self, req: &super::wallet::UnlockWalletRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::UnlockWalletResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_UNLOCK_WALLET, req, opt)
    }

    pub fn unlock_wallet_async(&self, req: &super::wallet::UnlockWalletRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::UnlockWalletResponse>> {
        self.unlock_wallet_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn sign_raw_transaction_opt(&self, req: &super::wallet::SignRawTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::SignRawTransactionResponse> {
        self.client.unary_call(&METHOD_WALLET_SIGN_RAW_TRANSACTION, req, opt)
    }

    pub fn sign_raw_transaction(&self, req: &super::wallet::SignRawTransactionRequest) -> ::grpcio::Result<super::wallet::SignRawTransactionResponse> {
        self.sign_raw_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn sign_raw_transaction_async_opt(&self, req: &super::wallet::SignRawTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::SignRawTransactionResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_SIGN_RAW_TRANSACTION, req, opt)
    }

    pub fn sign_raw_transaction_async(&self, req: &super::wallet::SignRawTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::SignRawTransactionResponse>> {
        self.sign_raw_transaction_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_validator_address_opt(&self, req: &super::wallet::GetValidatorAddressRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::wallet::GetValidatorAddressResponse> {
        self.client.unary_call(&METHOD_WALLET_GET_VALIDATOR_ADDRESS, req, opt)
    }

    pub fn get_validator_address(&self, req: &super::wallet::GetValidatorAddressRequest) -> ::grpcio::Result<super::wallet::GetValidatorAddressResponse> {
        self.get_validator_address_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_validator_address_async_opt(&self, req: &super::wallet::GetValidatorAddressRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::GetValidatorAddressResponse>> {
        self.client.unary_call_async(&METHOD_WALLET_GET_VALIDATOR_ADDRESS, req, opt)
    }

    pub fn get_validator_address_async(&self, req: &super::wallet::GetValidatorAddressRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::wallet::GetValidatorAddressResponse>> {
        self.get_validator_address_async_opt(req, ::grpcio::CallOption::default())
    }
    pub fn spawn<F>(&self, f: F) where F: ::std::future::Future<Output = ()> + Send + 'static {
        self.client.spawn(f)
    }
}

pub trait Wallet {
    fn create_wallet(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::CreateWalletRequest, sink: ::grpcio::UnarySink<super::wallet::CreateWalletResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn load_wallet(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::LoadWalletRequest, sink: ::grpcio::UnarySink<super::wallet::LoadWalletResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn unload_wallet(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::UnloadWalletRequest, sink: ::grpcio::UnarySink<super::wallet::UnloadWalletResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn lock_wallet(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::LockWalletRequest, sink: ::grpcio::UnarySink<super::wallet::LockWalletResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn unlock_wallet(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::UnlockWalletRequest, sink: ::grpcio::UnarySink<super::wallet::UnlockWalletResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn sign_raw_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::SignRawTransactionRequest, sink: ::grpcio::UnarySink<super::wallet::SignRawTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn get_validator_address(&mut self, ctx: ::grpcio::RpcContext, _req: super::wallet::GetValidatorAddressRequest, sink: ::grpcio::UnarySink<super::wallet::GetValidatorAddressResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
}

pub fn create_wallet<S: Wallet + Send + Clone + 'static>(s: S) -> ::grpcio::Service {
    let mut builder = ::grpcio::ServiceBuilder::new();
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_WALLET_CREATE_WALLET, move |ctx, req, resp| {
        instance.create_wallet(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_WALLET_LOAD_WALLET, move |ctx, req, resp| {
        instance.load_wallet(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_WALLET_UNLOAD_WALLET, move |ctx, req, resp| {
        instance.unload_wallet(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_WALLET_LOCK_WALLET, move |ctx, req, resp| {
        instance.lock_wallet(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_WALLET_UNLOCK_WALLET, move |ctx, req, resp| {
        instance.unlock_wallet(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_WALLET_SIGN_RAW_TRANSACTION, move |ctx, req, resp| {
        instance.sign_raw_transaction(ctx, req, resp)
    });
    let mut instance = s;
    builder = builder.add_unary_handler(&METHOD_WALLET_GET_VALIDATOR_ADDRESS, move |ctx, req, resp| {
        instance.get_validator_address(ctx, req, resp)
    });
    builder.build()
}
