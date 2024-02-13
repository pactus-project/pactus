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

const METHOD_TRANSACTION_GET_TRANSACTION: ::grpcio::Method<super::transaction::GetTransactionRequest, super::transaction::GetTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/GetTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_TRANSACTION_CALCULATE_FEE: ::grpcio::Method<super::transaction::CalculateFeeRequest, super::transaction::CalculateFeeResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/CalculateFee",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_TRANSACTION_BROADCAST_TRANSACTION: ::grpcio::Method<super::transaction::BroadcastTransactionRequest, super::transaction::BroadcastTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/BroadcastTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_TRANSACTION_GET_RAW_TRANSFER_TRANSACTION: ::grpcio::Method<super::transaction::GetRawTransferTransactionRequest, super::transaction::GetRawTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/GetRawTransferTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_TRANSACTION_GET_RAW_BOND_TRANSACTION: ::grpcio::Method<super::transaction::GetRawBondTransactionRequest, super::transaction::GetRawTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/GetRawBondTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_TRANSACTION_GET_RAW_UN_BOND_TRANSACTION: ::grpcio::Method<super::transaction::GetRawUnBondTransactionRequest, super::transaction::GetRawTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/GetRawUnBondTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_TRANSACTION_GET_RAW_WITHDRAW_TRANSACTION: ::grpcio::Method<super::transaction::GetRawWithdrawTransactionRequest, super::transaction::GetRawTransactionResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Transaction/GetRawWithdrawTransaction",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

#[derive(Clone)]
pub struct TransactionClient {
    pub client: ::grpcio::Client,
}

impl TransactionClient {
    pub fn new(channel: ::grpcio::Channel) -> Self {
        TransactionClient {
            client: ::grpcio::Client::new(channel),
        }
    }

    pub fn get_transaction_opt(&self, req: &super::transaction::GetTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::GetTransactionResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_GET_TRANSACTION, req, opt)
    }

    pub fn get_transaction(&self, req: &super::transaction::GetTransactionRequest) -> ::grpcio::Result<super::transaction::GetTransactionResponse> {
        self.get_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_transaction_async_opt(&self, req: &super::transaction::GetTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetTransactionResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_GET_TRANSACTION, req, opt)
    }

    pub fn get_transaction_async(&self, req: &super::transaction::GetTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetTransactionResponse>> {
        self.get_transaction_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn calculate_fee_opt(&self, req: &super::transaction::CalculateFeeRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::CalculateFeeResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_CALCULATE_FEE, req, opt)
    }

    pub fn calculate_fee(&self, req: &super::transaction::CalculateFeeRequest) -> ::grpcio::Result<super::transaction::CalculateFeeResponse> {
        self.calculate_fee_opt(req, ::grpcio::CallOption::default())
    }

    pub fn calculate_fee_async_opt(&self, req: &super::transaction::CalculateFeeRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::CalculateFeeResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_CALCULATE_FEE, req, opt)
    }

    pub fn calculate_fee_async(&self, req: &super::transaction::CalculateFeeRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::CalculateFeeResponse>> {
        self.calculate_fee_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn broadcast_transaction_opt(&self, req: &super::transaction::BroadcastTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::BroadcastTransactionResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_BROADCAST_TRANSACTION, req, opt)
    }

    pub fn broadcast_transaction(&self, req: &super::transaction::BroadcastTransactionRequest) -> ::grpcio::Result<super::transaction::BroadcastTransactionResponse> {
        self.broadcast_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn broadcast_transaction_async_opt(&self, req: &super::transaction::BroadcastTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::BroadcastTransactionResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_BROADCAST_TRANSACTION, req, opt)
    }

    pub fn broadcast_transaction_async(&self, req: &super::transaction::BroadcastTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::BroadcastTransactionResponse>> {
        self.broadcast_transaction_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_transfer_transaction_opt(&self, req: &super::transaction::GetRawTransferTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_GET_RAW_TRANSFER_TRANSACTION, req, opt)
    }

    pub fn get_raw_transfer_transaction(&self, req: &super::transaction::GetRawTransferTransactionRequest) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.get_raw_transfer_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_transfer_transaction_async_opt(&self, req: &super::transaction::GetRawTransferTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_GET_RAW_TRANSFER_TRANSACTION, req, opt)
    }

    pub fn get_raw_transfer_transaction_async(&self, req: &super::transaction::GetRawTransferTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.get_raw_transfer_transaction_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_bond_transaction_opt(&self, req: &super::transaction::GetRawBondTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_GET_RAW_BOND_TRANSACTION, req, opt)
    }

    pub fn get_raw_bond_transaction(&self, req: &super::transaction::GetRawBondTransactionRequest) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.get_raw_bond_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_bond_transaction_async_opt(&self, req: &super::transaction::GetRawBondTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_GET_RAW_BOND_TRANSACTION, req, opt)
    }

    pub fn get_raw_bond_transaction_async(&self, req: &super::transaction::GetRawBondTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.get_raw_bond_transaction_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_un_bond_transaction_opt(&self, req: &super::transaction::GetRawUnBondTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_GET_RAW_UN_BOND_TRANSACTION, req, opt)
    }

    pub fn get_raw_un_bond_transaction(&self, req: &super::transaction::GetRawUnBondTransactionRequest) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.get_raw_un_bond_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_un_bond_transaction_async_opt(&self, req: &super::transaction::GetRawUnBondTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_GET_RAW_UN_BOND_TRANSACTION, req, opt)
    }

    pub fn get_raw_un_bond_transaction_async(&self, req: &super::transaction::GetRawUnBondTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.get_raw_un_bond_transaction_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_withdraw_transaction_opt(&self, req: &super::transaction::GetRawWithdrawTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.client.unary_call(&METHOD_TRANSACTION_GET_RAW_WITHDRAW_TRANSACTION, req, opt)
    }

    pub fn get_raw_withdraw_transaction(&self, req: &super::transaction::GetRawWithdrawTransactionRequest) -> ::grpcio::Result<super::transaction::GetRawTransactionResponse> {
        self.get_raw_withdraw_transaction_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_raw_withdraw_transaction_async_opt(&self, req: &super::transaction::GetRawWithdrawTransactionRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.client.unary_call_async(&METHOD_TRANSACTION_GET_RAW_WITHDRAW_TRANSACTION, req, opt)
    }

    pub fn get_raw_withdraw_transaction_async(&self, req: &super::transaction::GetRawWithdrawTransactionRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::transaction::GetRawTransactionResponse>> {
        self.get_raw_withdraw_transaction_async_opt(req, ::grpcio::CallOption::default())
    }
    pub fn spawn<F>(&self, f: F) where F: ::std::future::Future<Output = ()> + Send + 'static {
        self.client.spawn(f)
    }
}

pub trait Transaction {
    fn get_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::GetTransactionRequest, sink: ::grpcio::UnarySink<super::transaction::GetTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn calculate_fee(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::CalculateFeeRequest, sink: ::grpcio::UnarySink<super::transaction::CalculateFeeResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn broadcast_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::BroadcastTransactionRequest, sink: ::grpcio::UnarySink<super::transaction::BroadcastTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn get_raw_transfer_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::GetRawTransferTransactionRequest, sink: ::grpcio::UnarySink<super::transaction::GetRawTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn get_raw_bond_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::GetRawBondTransactionRequest, sink: ::grpcio::UnarySink<super::transaction::GetRawTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn get_raw_un_bond_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::GetRawUnBondTransactionRequest, sink: ::grpcio::UnarySink<super::transaction::GetRawTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn get_raw_withdraw_transaction(&mut self, ctx: ::grpcio::RpcContext, _req: super::transaction::GetRawWithdrawTransactionRequest, sink: ::grpcio::UnarySink<super::transaction::GetRawTransactionResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
}

pub fn create_transaction<S: Transaction + Send + Clone + 'static>(s: S) -> ::grpcio::Service {
    let mut builder = ::grpcio::ServiceBuilder::new();
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_GET_TRANSACTION, move |ctx, req, resp| {
        instance.get_transaction(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_CALCULATE_FEE, move |ctx, req, resp| {
        instance.calculate_fee(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_BROADCAST_TRANSACTION, move |ctx, req, resp| {
        instance.broadcast_transaction(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_GET_RAW_TRANSFER_TRANSACTION, move |ctx, req, resp| {
        instance.get_raw_transfer_transaction(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_GET_RAW_BOND_TRANSACTION, move |ctx, req, resp| {
        instance.get_raw_bond_transaction(ctx, req, resp)
    });
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_GET_RAW_UN_BOND_TRANSACTION, move |ctx, req, resp| {
        instance.get_raw_un_bond_transaction(ctx, req, resp)
    });
    let mut instance = s;
    builder = builder.add_unary_handler(&METHOD_TRANSACTION_GET_RAW_WITHDRAW_TRANSACTION, move |ctx, req, resp| {
        instance.get_raw_withdraw_transaction(ctx, req, resp)
    });
    builder.build()
}
