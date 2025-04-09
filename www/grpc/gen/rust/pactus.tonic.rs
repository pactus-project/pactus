// @generated
/// Generated client implementations.
pub mod transaction_service_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct TransactionServiceClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl TransactionServiceClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> TransactionServiceClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> TransactionServiceClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + Send + Sync,
        {
            TransactionServiceClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn get_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::GetTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/GetTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.TransactionService", "GetTransaction"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn calculate_fee(
            &mut self,
            request: impl tonic::IntoRequest<super::CalculateFeeRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CalculateFeeResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/CalculateFee",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.TransactionService", "CalculateFee"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn broadcast_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::BroadcastTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::BroadcastTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/BroadcastTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.TransactionService", "BroadcastTransaction"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_raw_transfer_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::GetRawTransferTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/GetRawTransferTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "pactus.TransactionService",
                        "GetRawTransferTransaction",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_raw_bond_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::GetRawBondTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/GetRawBondTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.TransactionService", "GetRawBondTransaction"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_raw_unbond_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::GetRawUnbondTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/GetRawUnbondTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "pactus.TransactionService",
                        "GetRawUnbondTransaction",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_raw_withdraw_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::GetRawWithdrawTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/GetRawWithdrawTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "pactus.TransactionService",
                        "GetRawWithdrawTransaction",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn decode_raw_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::DecodeRawTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::DecodeRawTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.TransactionService/DecodeRawTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.TransactionService", "DecodeRawTransaction"),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod transaction_service_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with TransactionServiceServer.
    #[async_trait]
    pub trait TransactionService: Send + Sync + 'static {
        async fn get_transaction(
            &self,
            request: tonic::Request<super::GetTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTransactionResponse>,
            tonic::Status,
        >;
        async fn calculate_fee(
            &self,
            request: tonic::Request<super::CalculateFeeRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CalculateFeeResponse>,
            tonic::Status,
        >;
        async fn broadcast_transaction(
            &self,
            request: tonic::Request<super::BroadcastTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::BroadcastTransactionResponse>,
            tonic::Status,
        >;
        async fn get_raw_transfer_transaction(
            &self,
            request: tonic::Request<super::GetRawTransferTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        >;
        async fn get_raw_bond_transaction(
            &self,
            request: tonic::Request<super::GetRawBondTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        >;
        async fn get_raw_unbond_transaction(
            &self,
            request: tonic::Request<super::GetRawUnbondTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        >;
        async fn get_raw_withdraw_transaction(
            &self,
            request: tonic::Request<super::GetRawWithdrawTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetRawTransactionResponse>,
            tonic::Status,
        >;
        async fn decode_raw_transaction(
            &self,
            request: tonic::Request<super::DecodeRawTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::DecodeRawTransactionResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct TransactionServiceServer<T: TransactionService> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T: TransactionService> TransactionServiceServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for TransactionServiceServer<T>
    where
        T: TransactionService,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/pactus.TransactionService/GetTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct GetTransactionSvc<T: TransactionService>(pub Arc<T>);
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<super::GetTransactionRequest>
                    for GetTransactionSvc<T> {
                        type Response = super::GetTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetTransactionRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::get_transaction(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/CalculateFee" => {
                    #[allow(non_camel_case_types)]
                    struct CalculateFeeSvc<T: TransactionService>(pub Arc<T>);
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<super::CalculateFeeRequest>
                    for CalculateFeeSvc<T> {
                        type Response = super::CalculateFeeResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CalculateFeeRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::calculate_fee(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CalculateFeeSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/BroadcastTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct BroadcastTransactionSvc<T: TransactionService>(pub Arc<T>);
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<super::BroadcastTransactionRequest>
                    for BroadcastTransactionSvc<T> {
                        type Response = super::BroadcastTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::BroadcastTransactionRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::broadcast_transaction(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = BroadcastTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/GetRawTransferTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct GetRawTransferTransactionSvc<T: TransactionService>(
                        pub Arc<T>,
                    );
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<
                        super::GetRawTransferTransactionRequest,
                    > for GetRawTransferTransactionSvc<T> {
                        type Response = super::GetRawTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::GetRawTransferTransactionRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::get_raw_transfer_transaction(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetRawTransferTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/GetRawBondTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct GetRawBondTransactionSvc<T: TransactionService>(pub Arc<T>);
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<super::GetRawBondTransactionRequest>
                    for GetRawBondTransactionSvc<T> {
                        type Response = super::GetRawTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetRawBondTransactionRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::get_raw_bond_transaction(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetRawBondTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/GetRawUnbondTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct GetRawUnbondTransactionSvc<T: TransactionService>(pub Arc<T>);
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<super::GetRawUnbondTransactionRequest>
                    for GetRawUnbondTransactionSvc<T> {
                        type Response = super::GetRawTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::GetRawUnbondTransactionRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::get_raw_unbond_transaction(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetRawUnbondTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/GetRawWithdrawTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct GetRawWithdrawTransactionSvc<T: TransactionService>(
                        pub Arc<T>,
                    );
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<
                        super::GetRawWithdrawTransactionRequest,
                    > for GetRawWithdrawTransactionSvc<T> {
                        type Response = super::GetRawTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::GetRawWithdrawTransactionRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::get_raw_withdraw_transaction(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetRawWithdrawTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.TransactionService/DecodeRawTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct DecodeRawTransactionSvc<T: TransactionService>(pub Arc<T>);
                    impl<
                        T: TransactionService,
                    > tonic::server::UnaryService<super::DecodeRawTransactionRequest>
                    for DecodeRawTransactionSvc<T> {
                        type Response = super::DecodeRawTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::DecodeRawTransactionRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as TransactionService>::decode_raw_transaction(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = DecodeRawTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        Ok(
                            http::Response::builder()
                                .status(200)
                                .header("grpc-status", tonic::Code::Unimplemented as i32)
                                .header(
                                    http::header::CONTENT_TYPE,
                                    tonic::metadata::GRPC_CONTENT_TYPE,
                                )
                                .body(empty_body())
                                .unwrap(),
                        )
                    })
                }
            }
        }
    }
    impl<T: TransactionService> Clone for TransactionServiceServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    impl<T: TransactionService> tonic::server::NamedService
    for TransactionServiceServer<T> {
        const NAME: &'static str = "pactus.TransactionService";
    }
}
/// Generated client implementations.
pub mod blockchain_service_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct BlockchainServiceClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl BlockchainServiceClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> BlockchainServiceClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> BlockchainServiceClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + Send + Sync,
        {
            BlockchainServiceClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn get_block(
            &mut self,
            request: impl tonic::IntoRequest<super::GetBlockRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetBlock",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetBlock"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_block_hash(
            &mut self,
            request: impl tonic::IntoRequest<super::GetBlockHashRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockHashResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetBlockHash",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetBlockHash"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_block_height(
            &mut self,
            request: impl tonic::IntoRequest<super::GetBlockHeightRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockHeightResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetBlockHeight",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetBlockHeight"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_blockchain_info(
            &mut self,
            request: impl tonic::IntoRequest<super::GetBlockchainInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockchainInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetBlockchainInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.BlockchainService", "GetBlockchainInfo"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_consensus_info(
            &mut self,
            request: impl tonic::IntoRequest<super::GetConsensusInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetConsensusInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetConsensusInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetConsensusInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_account(
            &mut self,
            request: impl tonic::IntoRequest<super::GetAccountRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetAccountResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetAccount",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetAccount"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_validator(
            &mut self,
            request: impl tonic::IntoRequest<super::GetValidatorRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetValidator",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetValidator"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_validator_by_number(
            &mut self,
            request: impl tonic::IntoRequest<super::GetValidatorByNumberRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetValidatorByNumber",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.BlockchainService", "GetValidatorByNumber"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_validator_addresses(
            &mut self,
            request: impl tonic::IntoRequest<super::GetValidatorAddressesRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorAddressesResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetValidatorAddresses",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.BlockchainService", "GetValidatorAddresses"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_public_key(
            &mut self,
            request: impl tonic::IntoRequest<super::GetPublicKeyRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetPublicKeyResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetPublicKey",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetPublicKey"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_tx_pool_content(
            &mut self,
            request: impl tonic::IntoRequest<super::GetTxPoolContentRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTxPoolContentResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.BlockchainService/GetTxPoolContent",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.BlockchainService", "GetTxPoolContent"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod blockchain_service_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with BlockchainServiceServer.
    #[async_trait]
    pub trait BlockchainService: Send + Sync + 'static {
        async fn get_block(
            &self,
            request: tonic::Request<super::GetBlockRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockResponse>,
            tonic::Status,
        >;
        async fn get_block_hash(
            &self,
            request: tonic::Request<super::GetBlockHashRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockHashResponse>,
            tonic::Status,
        >;
        async fn get_block_height(
            &self,
            request: tonic::Request<super::GetBlockHeightRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockHeightResponse>,
            tonic::Status,
        >;
        async fn get_blockchain_info(
            &self,
            request: tonic::Request<super::GetBlockchainInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetBlockchainInfoResponse>,
            tonic::Status,
        >;
        async fn get_consensus_info(
            &self,
            request: tonic::Request<super::GetConsensusInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetConsensusInfoResponse>,
            tonic::Status,
        >;
        async fn get_account(
            &self,
            request: tonic::Request<super::GetAccountRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetAccountResponse>,
            tonic::Status,
        >;
        async fn get_validator(
            &self,
            request: tonic::Request<super::GetValidatorRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorResponse>,
            tonic::Status,
        >;
        async fn get_validator_by_number(
            &self,
            request: tonic::Request<super::GetValidatorByNumberRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorResponse>,
            tonic::Status,
        >;
        async fn get_validator_addresses(
            &self,
            request: tonic::Request<super::GetValidatorAddressesRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorAddressesResponse>,
            tonic::Status,
        >;
        async fn get_public_key(
            &self,
            request: tonic::Request<super::GetPublicKeyRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetPublicKeyResponse>,
            tonic::Status,
        >;
        async fn get_tx_pool_content(
            &self,
            request: tonic::Request<super::GetTxPoolContentRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTxPoolContentResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct BlockchainServiceServer<T: BlockchainService> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T: BlockchainService> BlockchainServiceServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for BlockchainServiceServer<T>
    where
        T: BlockchainService,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/pactus.BlockchainService/GetBlock" => {
                    #[allow(non_camel_case_types)]
                    struct GetBlockSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetBlockRequest>
                    for GetBlockSvc<T> {
                        type Response = super::GetBlockResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetBlockRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_block(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetBlockSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetBlockHash" => {
                    #[allow(non_camel_case_types)]
                    struct GetBlockHashSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetBlockHashRequest>
                    for GetBlockHashSvc<T> {
                        type Response = super::GetBlockHashResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetBlockHashRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_block_hash(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetBlockHashSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetBlockHeight" => {
                    #[allow(non_camel_case_types)]
                    struct GetBlockHeightSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetBlockHeightRequest>
                    for GetBlockHeightSvc<T> {
                        type Response = super::GetBlockHeightResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetBlockHeightRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_block_height(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetBlockHeightSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetBlockchainInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetBlockchainInfoSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetBlockchainInfoRequest>
                    for GetBlockchainInfoSvc<T> {
                        type Response = super::GetBlockchainInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetBlockchainInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_blockchain_info(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetBlockchainInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetConsensusInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetConsensusInfoSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetConsensusInfoRequest>
                    for GetConsensusInfoSvc<T> {
                        type Response = super::GetConsensusInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetConsensusInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_consensus_info(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetConsensusInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetAccount" => {
                    #[allow(non_camel_case_types)]
                    struct GetAccountSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetAccountRequest>
                    for GetAccountSvc<T> {
                        type Response = super::GetAccountResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetAccountRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_account(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetAccountSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetValidator" => {
                    #[allow(non_camel_case_types)]
                    struct GetValidatorSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetValidatorRequest>
                    for GetValidatorSvc<T> {
                        type Response = super::GetValidatorResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetValidatorRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_validator(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetValidatorSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetValidatorByNumber" => {
                    #[allow(non_camel_case_types)]
                    struct GetValidatorByNumberSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetValidatorByNumberRequest>
                    for GetValidatorByNumberSvc<T> {
                        type Response = super::GetValidatorResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetValidatorByNumberRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_validator_by_number(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetValidatorByNumberSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetValidatorAddresses" => {
                    #[allow(non_camel_case_types)]
                    struct GetValidatorAddressesSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetValidatorAddressesRequest>
                    for GetValidatorAddressesSvc<T> {
                        type Response = super::GetValidatorAddressesResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetValidatorAddressesRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_validator_addresses(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetValidatorAddressesSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetPublicKey" => {
                    #[allow(non_camel_case_types)]
                    struct GetPublicKeySvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetPublicKeyRequest>
                    for GetPublicKeySvc<T> {
                        type Response = super::GetPublicKeyResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetPublicKeyRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_public_key(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetPublicKeySvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.BlockchainService/GetTxPoolContent" => {
                    #[allow(non_camel_case_types)]
                    struct GetTxPoolContentSvc<T: BlockchainService>(pub Arc<T>);
                    impl<
                        T: BlockchainService,
                    > tonic::server::UnaryService<super::GetTxPoolContentRequest>
                    for GetTxPoolContentSvc<T> {
                        type Response = super::GetTxPoolContentResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetTxPoolContentRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as BlockchainService>::get_tx_pool_content(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetTxPoolContentSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        Ok(
                            http::Response::builder()
                                .status(200)
                                .header("grpc-status", tonic::Code::Unimplemented as i32)
                                .header(
                                    http::header::CONTENT_TYPE,
                                    tonic::metadata::GRPC_CONTENT_TYPE,
                                )
                                .body(empty_body())
                                .unwrap(),
                        )
                    })
                }
            }
        }
    }
    impl<T: BlockchainService> Clone for BlockchainServiceServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    impl<T: BlockchainService> tonic::server::NamedService
    for BlockchainServiceServer<T> {
        const NAME: &'static str = "pactus.BlockchainService";
    }
}
/// Generated client implementations.
pub mod network_service_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct NetworkServiceClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl NetworkServiceClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> NetworkServiceClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> NetworkServiceClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + Send + Sync,
        {
            NetworkServiceClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn get_network_info(
            &mut self,
            request: impl tonic::IntoRequest<super::GetNetworkInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetNetworkInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.NetworkService/GetNetworkInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.NetworkService", "GetNetworkInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_node_info(
            &mut self,
            request: impl tonic::IntoRequest<super::GetNodeInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetNodeInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.NetworkService/GetNodeInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.NetworkService", "GetNodeInfo"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod network_service_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with NetworkServiceServer.
    #[async_trait]
    pub trait NetworkService: Send + Sync + 'static {
        async fn get_network_info(
            &self,
            request: tonic::Request<super::GetNetworkInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetNetworkInfoResponse>,
            tonic::Status,
        >;
        async fn get_node_info(
            &self,
            request: tonic::Request<super::GetNodeInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetNodeInfoResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct NetworkServiceServer<T: NetworkService> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T: NetworkService> NetworkServiceServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for NetworkServiceServer<T>
    where
        T: NetworkService,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/pactus.NetworkService/GetNetworkInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetNetworkInfoSvc<T: NetworkService>(pub Arc<T>);
                    impl<
                        T: NetworkService,
                    > tonic::server::UnaryService<super::GetNetworkInfoRequest>
                    for GetNetworkInfoSvc<T> {
                        type Response = super::GetNetworkInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetNetworkInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as NetworkService>::get_network_info(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetNetworkInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.NetworkService/GetNodeInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetNodeInfoSvc<T: NetworkService>(pub Arc<T>);
                    impl<
                        T: NetworkService,
                    > tonic::server::UnaryService<super::GetNodeInfoRequest>
                    for GetNodeInfoSvc<T> {
                        type Response = super::GetNodeInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetNodeInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as NetworkService>::get_node_info(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetNodeInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        Ok(
                            http::Response::builder()
                                .status(200)
                                .header("grpc-status", tonic::Code::Unimplemented as i32)
                                .header(
                                    http::header::CONTENT_TYPE,
                                    tonic::metadata::GRPC_CONTENT_TYPE,
                                )
                                .body(empty_body())
                                .unwrap(),
                        )
                    })
                }
            }
        }
    }
    impl<T: NetworkService> Clone for NetworkServiceServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    impl<T: NetworkService> tonic::server::NamedService for NetworkServiceServer<T> {
        const NAME: &'static str = "pactus.NetworkService";
    }
}
/// Generated client implementations.
pub mod utils_service_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct UtilsServiceClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl UtilsServiceClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> UtilsServiceClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> UtilsServiceClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + Send + Sync,
        {
            UtilsServiceClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn sign_message_with_private_key(
            &mut self,
            request: impl tonic::IntoRequest<super::SignMessageWithPrivateKeyRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignMessageWithPrivateKeyResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.UtilsService/SignMessageWithPrivateKey",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("pactus.UtilsService", "SignMessageWithPrivateKey"),
                );
            self.inner.unary(req, path, codec).await
        }
        pub async fn verify_message(
            &mut self,
            request: impl tonic::IntoRequest<super::VerifyMessageRequest>,
        ) -> std::result::Result<
            tonic::Response<super::VerifyMessageResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.UtilsService/VerifyMessage",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.UtilsService", "VerifyMessage"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn public_key_aggregation(
            &mut self,
            request: impl tonic::IntoRequest<super::PublicKeyAggregationRequest>,
        ) -> std::result::Result<
            tonic::Response<super::PublicKeyAggregationResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.UtilsService/PublicKeyAggregation",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.UtilsService", "PublicKeyAggregation"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn signature_aggregation(
            &mut self,
            request: impl tonic::IntoRequest<super::SignatureAggregationRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignatureAggregationResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.UtilsService/SignatureAggregation",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.UtilsService", "SignatureAggregation"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod utils_service_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UtilsServiceServer.
    #[async_trait]
    pub trait UtilsService: Send + Sync + 'static {
        async fn sign_message_with_private_key(
            &self,
            request: tonic::Request<super::SignMessageWithPrivateKeyRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignMessageWithPrivateKeyResponse>,
            tonic::Status,
        >;
        async fn verify_message(
            &self,
            request: tonic::Request<super::VerifyMessageRequest>,
        ) -> std::result::Result<
            tonic::Response<super::VerifyMessageResponse>,
            tonic::Status,
        >;
        async fn public_key_aggregation(
            &self,
            request: tonic::Request<super::PublicKeyAggregationRequest>,
        ) -> std::result::Result<
            tonic::Response<super::PublicKeyAggregationResponse>,
            tonic::Status,
        >;
        async fn signature_aggregation(
            &self,
            request: tonic::Request<super::SignatureAggregationRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignatureAggregationResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct UtilsServiceServer<T: UtilsService> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T: UtilsService> UtilsServiceServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UtilsServiceServer<T>
    where
        T: UtilsService,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/pactus.UtilsService/SignMessageWithPrivateKey" => {
                    #[allow(non_camel_case_types)]
                    struct SignMessageWithPrivateKeySvc<T: UtilsService>(pub Arc<T>);
                    impl<
                        T: UtilsService,
                    > tonic::server::UnaryService<
                        super::SignMessageWithPrivateKeyRequest,
                    > for SignMessageWithPrivateKeySvc<T> {
                        type Response = super::SignMessageWithPrivateKeyResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::SignMessageWithPrivateKeyRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as UtilsService>::sign_message_with_private_key(
                                        &inner,
                                        request,
                                    )
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SignMessageWithPrivateKeySvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.UtilsService/VerifyMessage" => {
                    #[allow(non_camel_case_types)]
                    struct VerifyMessageSvc<T: UtilsService>(pub Arc<T>);
                    impl<
                        T: UtilsService,
                    > tonic::server::UnaryService<super::VerifyMessageRequest>
                    for VerifyMessageSvc<T> {
                        type Response = super::VerifyMessageResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::VerifyMessageRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as UtilsService>::verify_message(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = VerifyMessageSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.UtilsService/PublicKeyAggregation" => {
                    #[allow(non_camel_case_types)]
                    struct PublicKeyAggregationSvc<T: UtilsService>(pub Arc<T>);
                    impl<
                        T: UtilsService,
                    > tonic::server::UnaryService<super::PublicKeyAggregationRequest>
                    for PublicKeyAggregationSvc<T> {
                        type Response = super::PublicKeyAggregationResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::PublicKeyAggregationRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as UtilsService>::public_key_aggregation(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = PublicKeyAggregationSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.UtilsService/SignatureAggregation" => {
                    #[allow(non_camel_case_types)]
                    struct SignatureAggregationSvc<T: UtilsService>(pub Arc<T>);
                    impl<
                        T: UtilsService,
                    > tonic::server::UnaryService<super::SignatureAggregationRequest>
                    for SignatureAggregationSvc<T> {
                        type Response = super::SignatureAggregationResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SignatureAggregationRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as UtilsService>::signature_aggregation(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SignatureAggregationSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        Ok(
                            http::Response::builder()
                                .status(200)
                                .header("grpc-status", tonic::Code::Unimplemented as i32)
                                .header(
                                    http::header::CONTENT_TYPE,
                                    tonic::metadata::GRPC_CONTENT_TYPE,
                                )
                                .body(empty_body())
                                .unwrap(),
                        )
                    })
                }
            }
        }
    }
    impl<T: UtilsService> Clone for UtilsServiceServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    impl<T: UtilsService> tonic::server::NamedService for UtilsServiceServer<T> {
        const NAME: &'static str = "pactus.UtilsService";
    }
}
/// Generated client implementations.
pub mod wallet_service_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct WalletServiceClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl WalletServiceClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> WalletServiceClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> WalletServiceClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + Send + Sync,
        {
            WalletServiceClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn create_wallet(
            &mut self,
            request: impl tonic::IntoRequest<super::CreateWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CreateWalletResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/CreateWallet",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "CreateWallet"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn restore_wallet(
            &mut self,
            request: impl tonic::IntoRequest<super::RestoreWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::RestoreWalletResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/RestoreWallet",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "RestoreWallet"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn load_wallet(
            &mut self,
            request: impl tonic::IntoRequest<super::LoadWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::LoadWalletResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/LoadWallet",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "LoadWallet"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn unload_wallet(
            &mut self,
            request: impl tonic::IntoRequest<super::UnloadWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::UnloadWalletResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/UnloadWallet",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "UnloadWallet"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_total_balance(
            &mut self,
            request: impl tonic::IntoRequest<super::GetTotalBalanceRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTotalBalanceResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetTotalBalance",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetTotalBalance"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn sign_raw_transaction(
            &mut self,
            request: impl tonic::IntoRequest<super::SignRawTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignRawTransactionResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/SignRawTransaction",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "SignRawTransaction"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_validator_address(
            &mut self,
            request: impl tonic::IntoRequest<super::GetValidatorAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorAddressResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetValidatorAddress",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetValidatorAddress"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_new_address(
            &mut self,
            request: impl tonic::IntoRequest<super::GetNewAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetNewAddressResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetNewAddress",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetNewAddress"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_address_history(
            &mut self,
            request: impl tonic::IntoRequest<super::GetAddressHistoryRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetAddressHistoryResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetAddressHistory",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetAddressHistory"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn sign_message(
            &mut self,
            request: impl tonic::IntoRequest<super::SignMessageRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignMessageResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/SignMessage",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "SignMessage"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_total_stake(
            &mut self,
            request: impl tonic::IntoRequest<super::GetTotalStakeRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTotalStakeResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetTotalStake",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetTotalStake"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_address_info(
            &mut self,
            request: impl tonic::IntoRequest<super::GetAddressInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetAddressInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetAddressInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetAddressInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn set_address_label(
            &mut self,
            request: impl tonic::IntoRequest<super::SetAddressLabelRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SetAddressLabelResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/SetAddressLabel",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "SetAddressLabel"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn list_wallet(
            &mut self,
            request: impl tonic::IntoRequest<super::ListWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::ListWalletResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/ListWallet",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "ListWallet"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_wallet_info(
            &mut self,
            request: impl tonic::IntoRequest<super::GetWalletInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetWalletInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/GetWalletInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "GetWalletInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn list_address(
            &mut self,
            request: impl tonic::IntoRequest<super::ListAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::ListAddressResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/pactus.WalletService/ListAddress",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("pactus.WalletService", "ListAddress"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod wallet_service_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with WalletServiceServer.
    #[async_trait]
    pub trait WalletService: Send + Sync + 'static {
        async fn create_wallet(
            &self,
            request: tonic::Request<super::CreateWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::CreateWalletResponse>,
            tonic::Status,
        >;
        async fn restore_wallet(
            &self,
            request: tonic::Request<super::RestoreWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::RestoreWalletResponse>,
            tonic::Status,
        >;
        async fn load_wallet(
            &self,
            request: tonic::Request<super::LoadWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::LoadWalletResponse>,
            tonic::Status,
        >;
        async fn unload_wallet(
            &self,
            request: tonic::Request<super::UnloadWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::UnloadWalletResponse>,
            tonic::Status,
        >;
        async fn get_total_balance(
            &self,
            request: tonic::Request<super::GetTotalBalanceRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTotalBalanceResponse>,
            tonic::Status,
        >;
        async fn sign_raw_transaction(
            &self,
            request: tonic::Request<super::SignRawTransactionRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignRawTransactionResponse>,
            tonic::Status,
        >;
        async fn get_validator_address(
            &self,
            request: tonic::Request<super::GetValidatorAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetValidatorAddressResponse>,
            tonic::Status,
        >;
        async fn get_new_address(
            &self,
            request: tonic::Request<super::GetNewAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetNewAddressResponse>,
            tonic::Status,
        >;
        async fn get_address_history(
            &self,
            request: tonic::Request<super::GetAddressHistoryRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetAddressHistoryResponse>,
            tonic::Status,
        >;
        async fn sign_message(
            &self,
            request: tonic::Request<super::SignMessageRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SignMessageResponse>,
            tonic::Status,
        >;
        async fn get_total_stake(
            &self,
            request: tonic::Request<super::GetTotalStakeRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetTotalStakeResponse>,
            tonic::Status,
        >;
        async fn get_address_info(
            &self,
            request: tonic::Request<super::GetAddressInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetAddressInfoResponse>,
            tonic::Status,
        >;
        async fn set_address_label(
            &self,
            request: tonic::Request<super::SetAddressLabelRequest>,
        ) -> std::result::Result<
            tonic::Response<super::SetAddressLabelResponse>,
            tonic::Status,
        >;
        async fn list_wallet(
            &self,
            request: tonic::Request<super::ListWalletRequest>,
        ) -> std::result::Result<
            tonic::Response<super::ListWalletResponse>,
            tonic::Status,
        >;
        async fn get_wallet_info(
            &self,
            request: tonic::Request<super::GetWalletInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::GetWalletInfoResponse>,
            tonic::Status,
        >;
        async fn list_address(
            &self,
            request: tonic::Request<super::ListAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::ListAddressResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct WalletServiceServer<T: WalletService> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T: WalletService> WalletServiceServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for WalletServiceServer<T>
    where
        T: WalletService,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/pactus.WalletService/CreateWallet" => {
                    #[allow(non_camel_case_types)]
                    struct CreateWalletSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::CreateWalletRequest>
                    for CreateWalletSvc<T> {
                        type Response = super::CreateWalletResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::CreateWalletRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::create_wallet(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = CreateWalletSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/RestoreWallet" => {
                    #[allow(non_camel_case_types)]
                    struct RestoreWalletSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::RestoreWalletRequest>
                    for RestoreWalletSvc<T> {
                        type Response = super::RestoreWalletResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::RestoreWalletRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::restore_wallet(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RestoreWalletSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/LoadWallet" => {
                    #[allow(non_camel_case_types)]
                    struct LoadWalletSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::LoadWalletRequest>
                    for LoadWalletSvc<T> {
                        type Response = super::LoadWalletResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::LoadWalletRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::load_wallet(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = LoadWalletSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/UnloadWallet" => {
                    #[allow(non_camel_case_types)]
                    struct UnloadWalletSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::UnloadWalletRequest>
                    for UnloadWalletSvc<T> {
                        type Response = super::UnloadWalletResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::UnloadWalletRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::unload_wallet(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UnloadWalletSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetTotalBalance" => {
                    #[allow(non_camel_case_types)]
                    struct GetTotalBalanceSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetTotalBalanceRequest>
                    for GetTotalBalanceSvc<T> {
                        type Response = super::GetTotalBalanceResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetTotalBalanceRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_total_balance(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetTotalBalanceSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/SignRawTransaction" => {
                    #[allow(non_camel_case_types)]
                    struct SignRawTransactionSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::SignRawTransactionRequest>
                    for SignRawTransactionSvc<T> {
                        type Response = super::SignRawTransactionResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SignRawTransactionRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::sign_raw_transaction(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SignRawTransactionSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetValidatorAddress" => {
                    #[allow(non_camel_case_types)]
                    struct GetValidatorAddressSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetValidatorAddressRequest>
                    for GetValidatorAddressSvc<T> {
                        type Response = super::GetValidatorAddressResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetValidatorAddressRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_validator_address(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetValidatorAddressSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetNewAddress" => {
                    #[allow(non_camel_case_types)]
                    struct GetNewAddressSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetNewAddressRequest>
                    for GetNewAddressSvc<T> {
                        type Response = super::GetNewAddressResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetNewAddressRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_new_address(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetNewAddressSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetAddressHistory" => {
                    #[allow(non_camel_case_types)]
                    struct GetAddressHistorySvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetAddressHistoryRequest>
                    for GetAddressHistorySvc<T> {
                        type Response = super::GetAddressHistoryResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetAddressHistoryRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_address_history(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetAddressHistorySvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/SignMessage" => {
                    #[allow(non_camel_case_types)]
                    struct SignMessageSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::SignMessageRequest>
                    for SignMessageSvc<T> {
                        type Response = super::SignMessageResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SignMessageRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::sign_message(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SignMessageSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetTotalStake" => {
                    #[allow(non_camel_case_types)]
                    struct GetTotalStakeSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetTotalStakeRequest>
                    for GetTotalStakeSvc<T> {
                        type Response = super::GetTotalStakeResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetTotalStakeRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_total_stake(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetTotalStakeSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetAddressInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetAddressInfoSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetAddressInfoRequest>
                    for GetAddressInfoSvc<T> {
                        type Response = super::GetAddressInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetAddressInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_address_info(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetAddressInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/SetAddressLabel" => {
                    #[allow(non_camel_case_types)]
                    struct SetAddressLabelSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::SetAddressLabelRequest>
                    for SetAddressLabelSvc<T> {
                        type Response = super::SetAddressLabelResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::SetAddressLabelRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::set_address_label(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SetAddressLabelSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/ListWallet" => {
                    #[allow(non_camel_case_types)]
                    struct ListWalletSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::ListWalletRequest>
                    for ListWalletSvc<T> {
                        type Response = super::ListWalletResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ListWalletRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::list_wallet(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = ListWalletSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/GetWalletInfo" => {
                    #[allow(non_camel_case_types)]
                    struct GetWalletInfoSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::GetWalletInfoRequest>
                    for GetWalletInfoSvc<T> {
                        type Response = super::GetWalletInfoResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::GetWalletInfoRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::get_wallet_info(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetWalletInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/pactus.WalletService/ListAddress" => {
                    #[allow(non_camel_case_types)]
                    struct ListAddressSvc<T: WalletService>(pub Arc<T>);
                    impl<
                        T: WalletService,
                    > tonic::server::UnaryService<super::ListAddressRequest>
                    for ListAddressSvc<T> {
                        type Response = super::ListAddressResponse;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ListAddressRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as WalletService>::list_address(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = ListAddressSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        Ok(
                            http::Response::builder()
                                .status(200)
                                .header("grpc-status", tonic::Code::Unimplemented as i32)
                                .header(
                                    http::header::CONTENT_TYPE,
                                    tonic::metadata::GRPC_CONTENT_TYPE,
                                )
                                .body(empty_body())
                                .unwrap(),
                        )
                    })
                }
            }
        }
    }
    impl<T: WalletService> Clone for WalletServiceServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    impl<T: WalletService> tonic::server::NamedService for WalletServiceServer<T> {
        const NAME: &'static str = "pactus.WalletService";
    }
}
