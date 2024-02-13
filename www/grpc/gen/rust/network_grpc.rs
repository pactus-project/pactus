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

const METHOD_NETWORK_GET_NETWORK_INFO: ::grpcio::Method<super::network::GetNetworkInfoRequest, super::network::GetNetworkInfoResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Network/GetNetworkInfo",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

const METHOD_NETWORK_GET_NODE_INFO: ::grpcio::Method<super::network::GetNodeInfoRequest, super::network::GetNodeInfoResponse> = ::grpcio::Method {
    ty: ::grpcio::MethodType::Unary,
    name: "/pactus.Network/GetNodeInfo",
    req_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
    resp_mar: ::grpcio::Marshaller { ser: ::grpcio::pb_ser, de: ::grpcio::pb_de },
};

#[derive(Clone)]
pub struct NetworkClient {
    pub client: ::grpcio::Client,
}

impl NetworkClient {
    pub fn new(channel: ::grpcio::Channel) -> Self {
        NetworkClient {
            client: ::grpcio::Client::new(channel),
        }
    }

    pub fn get_network_info_opt(&self, req: &super::network::GetNetworkInfoRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::network::GetNetworkInfoResponse> {
        self.client.unary_call(&METHOD_NETWORK_GET_NETWORK_INFO, req, opt)
    }

    pub fn get_network_info(&self, req: &super::network::GetNetworkInfoRequest) -> ::grpcio::Result<super::network::GetNetworkInfoResponse> {
        self.get_network_info_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_network_info_async_opt(&self, req: &super::network::GetNetworkInfoRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::network::GetNetworkInfoResponse>> {
        self.client.unary_call_async(&METHOD_NETWORK_GET_NETWORK_INFO, req, opt)
    }

    pub fn get_network_info_async(&self, req: &super::network::GetNetworkInfoRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::network::GetNetworkInfoResponse>> {
        self.get_network_info_async_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_node_info_opt(&self, req: &super::network::GetNodeInfoRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<super::network::GetNodeInfoResponse> {
        self.client.unary_call(&METHOD_NETWORK_GET_NODE_INFO, req, opt)
    }

    pub fn get_node_info(&self, req: &super::network::GetNodeInfoRequest) -> ::grpcio::Result<super::network::GetNodeInfoResponse> {
        self.get_node_info_opt(req, ::grpcio::CallOption::default())
    }

    pub fn get_node_info_async_opt(&self, req: &super::network::GetNodeInfoRequest, opt: ::grpcio::CallOption) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::network::GetNodeInfoResponse>> {
        self.client.unary_call_async(&METHOD_NETWORK_GET_NODE_INFO, req, opt)
    }

    pub fn get_node_info_async(&self, req: &super::network::GetNodeInfoRequest) -> ::grpcio::Result<::grpcio::ClientUnaryReceiver<super::network::GetNodeInfoResponse>> {
        self.get_node_info_async_opt(req, ::grpcio::CallOption::default())
    }
    pub fn spawn<F>(&self, f: F) where F: ::std::future::Future<Output = ()> + Send + 'static {
        self.client.spawn(f)
    }
}

pub trait Network {
    fn get_network_info(&mut self, ctx: ::grpcio::RpcContext, _req: super::network::GetNetworkInfoRequest, sink: ::grpcio::UnarySink<super::network::GetNetworkInfoResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
    fn get_node_info(&mut self, ctx: ::grpcio::RpcContext, _req: super::network::GetNodeInfoRequest, sink: ::grpcio::UnarySink<super::network::GetNodeInfoResponse>) {
        grpcio::unimplemented_call!(ctx, sink)
    }
}

pub fn create_network<S: Network + Send + Clone + 'static>(s: S) -> ::grpcio::Service {
    let mut builder = ::grpcio::ServiceBuilder::new();
    let mut instance = s.clone();
    builder = builder.add_unary_handler(&METHOD_NETWORK_GET_NETWORK_INFO, move |ctx, req, resp| {
        instance.get_network_info(ctx, req, resp)
    });
    let mut instance = s;
    builder = builder.add_unary_handler(&METHOD_NETWORK_GET_NODE_INFO, move |ctx, req, resp| {
        instance.get_node_info(ctx, req, resp)
    });
    builder.build()
}
