// Code generated by MockGen. DO NOT EDIT.
// Source: libp2p.go

// Package mock_libp2p is a generated GoMock package.
package mock_libp2p

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	goprocess "github.com/jbenet/goprocess"
	connmgr "github.com/libp2p/go-libp2p-core/connmgr"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	event "github.com/libp2p/go-libp2p-core/event"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// MockHost is a mock of Host interface.
type MockHost struct {
	ctrl     *gomock.Controller
	recorder *MockHostMockRecorder
}

// MockHostMockRecorder is the mock recorder for MockHost.
type MockHostMockRecorder struct {
	mock *MockHost
}

// NewMockHost creates a new mock instance.
func NewMockHost(ctrl *gomock.Controller) *MockHost {
	mock := &MockHost{ctrl: ctrl}
	mock.recorder = &MockHostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHost) EXPECT() *MockHostMockRecorder {
	return m.recorder
}

// Addrs mocks base method.
func (m *MockHost) Addrs() []multiaddr.Multiaddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addrs")
	ret0, _ := ret[0].([]multiaddr.Multiaddr)
	return ret0
}

// Addrs indicates an expected call of Addrs.
func (mr *MockHostMockRecorder) Addrs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addrs", reflect.TypeOf((*MockHost)(nil).Addrs))
}

// Close mocks base method.
func (m *MockHost) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockHostMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockHost)(nil).Close))
}

// ConnManager mocks base method.
func (m *MockHost) ConnManager() connmgr.ConnManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnManager")
	ret0, _ := ret[0].(connmgr.ConnManager)
	return ret0
}

// ConnManager indicates an expected call of ConnManager.
func (mr *MockHostMockRecorder) ConnManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnManager", reflect.TypeOf((*MockHost)(nil).ConnManager))
}

// Connect mocks base method.
func (m *MockHost) Connect(ctx context.Context, pi peer.AddrInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", ctx, pi)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockHostMockRecorder) Connect(ctx, pi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockHost)(nil).Connect), ctx, pi)
}

// EventBus mocks base method.
func (m *MockHost) EventBus() event.Bus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventBus")
	ret0, _ := ret[0].(event.Bus)
	return ret0
}

// EventBus indicates an expected call of EventBus.
func (mr *MockHostMockRecorder) EventBus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventBus", reflect.TypeOf((*MockHost)(nil).EventBus))
}

// ID mocks base method.
func (m *MockHost) ID() peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(peer.ID)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockHostMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockHost)(nil).ID))
}

// Mux mocks base method.
func (m *MockHost) Mux() protocol.Switch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mux")
	ret0, _ := ret[0].(protocol.Switch)
	return ret0
}

// Mux indicates an expected call of Mux.
func (mr *MockHostMockRecorder) Mux() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mux", reflect.TypeOf((*MockHost)(nil).Mux))
}

// Network mocks base method.
func (m *MockHost) Network() network.Network {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Network")
	ret0, _ := ret[0].(network.Network)
	return ret0
}

// Network indicates an expected call of Network.
func (mr *MockHostMockRecorder) Network() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Network", reflect.TypeOf((*MockHost)(nil).Network))
}

// NewStream mocks base method.
func (m *MockHost) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, p}
	for _, a := range pids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewStream", varargs...)
	ret0, _ := ret[0].(network.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream.
func (mr *MockHostMockRecorder) NewStream(ctx, p interface{}, pids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, p}, pids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*MockHost)(nil).NewStream), varargs...)
}

// Peerstore mocks base method.
func (m *MockHost) Peerstore() peerstore.Peerstore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peerstore")
	ret0, _ := ret[0].(peerstore.Peerstore)
	return ret0
}

// Peerstore indicates an expected call of Peerstore.
func (mr *MockHostMockRecorder) Peerstore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peerstore", reflect.TypeOf((*MockHost)(nil).Peerstore))
}

// RemoveStreamHandler mocks base method.
func (m *MockHost) RemoveStreamHandler(pid protocol.ID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveStreamHandler", pid)
}

// RemoveStreamHandler indicates an expected call of RemoveStreamHandler.
func (mr *MockHostMockRecorder) RemoveStreamHandler(pid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveStreamHandler", reflect.TypeOf((*MockHost)(nil).RemoveStreamHandler), pid)
}

// SetStreamHandler mocks base method.
func (m *MockHost) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStreamHandler", pid, handler)
}

// SetStreamHandler indicates an expected call of SetStreamHandler.
func (mr *MockHostMockRecorder) SetStreamHandler(pid, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStreamHandler", reflect.TypeOf((*MockHost)(nil).SetStreamHandler), pid, handler)
}

// SetStreamHandlerMatch mocks base method.
func (m *MockHost) SetStreamHandlerMatch(arg0 protocol.ID, arg1 func(string) bool, arg2 network.StreamHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStreamHandlerMatch", arg0, arg1, arg2)
}

// SetStreamHandlerMatch indicates an expected call of SetStreamHandlerMatch.
func (mr *MockHostMockRecorder) SetStreamHandlerMatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStreamHandlerMatch", reflect.TypeOf((*MockHost)(nil).SetStreamHandlerMatch), arg0, arg1, arg2)
}

// MockConn is a mock of Conn interface.
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn.
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance.
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
}

// GetStreams mocks base method.
func (m *MockConn) GetStreams() []network.Stream {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStreams")
	ret0, _ := ret[0].([]network.Stream)
	return ret0
}

// GetStreams indicates an expected call of GetStreams.
func (mr *MockConnMockRecorder) GetStreams() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStreams", reflect.TypeOf((*MockConn)(nil).GetStreams))
}

// ID mocks base method.
func (m *MockConn) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockConnMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockConn)(nil).ID))
}

// LocalMultiaddr mocks base method.
func (m *MockConn) LocalMultiaddr() multiaddr.Multiaddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalMultiaddr")
	ret0, _ := ret[0].(multiaddr.Multiaddr)
	return ret0
}

// LocalMultiaddr indicates an expected call of LocalMultiaddr.
func (mr *MockConnMockRecorder) LocalMultiaddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalMultiaddr", reflect.TypeOf((*MockConn)(nil).LocalMultiaddr))
}

// LocalPeer mocks base method.
func (m *MockConn) LocalPeer() peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalPeer")
	ret0, _ := ret[0].(peer.ID)
	return ret0
}

// LocalPeer indicates an expected call of LocalPeer.
func (mr *MockConnMockRecorder) LocalPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalPeer", reflect.TypeOf((*MockConn)(nil).LocalPeer))
}

// LocalPrivateKey mocks base method.
func (m *MockConn) LocalPrivateKey() crypto.PrivKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalPrivateKey")
	ret0, _ := ret[0].(crypto.PrivKey)
	return ret0
}

// LocalPrivateKey indicates an expected call of LocalPrivateKey.
func (mr *MockConnMockRecorder) LocalPrivateKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalPrivateKey", reflect.TypeOf((*MockConn)(nil).LocalPrivateKey))
}

// NewStream mocks base method.
func (m *MockConn) NewStream(arg0 context.Context) (network.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewStream", arg0)
	ret0, _ := ret[0].(network.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream.
func (mr *MockConnMockRecorder) NewStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*MockConn)(nil).NewStream), arg0)
}

// RemoteMultiaddr mocks base method.
func (m *MockConn) RemoteMultiaddr() multiaddr.Multiaddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteMultiaddr")
	ret0, _ := ret[0].(multiaddr.Multiaddr)
	return ret0
}

// RemoteMultiaddr indicates an expected call of RemoteMultiaddr.
func (mr *MockConnMockRecorder) RemoteMultiaddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteMultiaddr", reflect.TypeOf((*MockConn)(nil).RemoteMultiaddr))
}

// RemotePeer mocks base method.
func (m *MockConn) RemotePeer() peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemotePeer")
	ret0, _ := ret[0].(peer.ID)
	return ret0
}

// RemotePeer indicates an expected call of RemotePeer.
func (mr *MockConnMockRecorder) RemotePeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemotePeer", reflect.TypeOf((*MockConn)(nil).RemotePeer))
}

// RemotePublicKey mocks base method.
func (m *MockConn) RemotePublicKey() crypto.PubKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemotePublicKey")
	ret0, _ := ret[0].(crypto.PubKey)
	return ret0
}

// RemotePublicKey indicates an expected call of RemotePublicKey.
func (mr *MockConnMockRecorder) RemotePublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemotePublicKey", reflect.TypeOf((*MockConn)(nil).RemotePublicKey))
}

// Stat mocks base method.
func (m *MockConn) Stat() network.Stat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(network.Stat)
	return ret0
}

// Stat indicates an expected call of Stat.
func (mr *MockConnMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockConn)(nil).Stat))
}

// MockStream is a mock of Stream interface.
type MockStream struct {
	ctrl     *gomock.Controller
	recorder *MockStreamMockRecorder
}

// MockStreamMockRecorder is the mock recorder for MockStream.
type MockStreamMockRecorder struct {
	mock *MockStream
}

// NewMockStream creates a new mock instance.
func NewMockStream(ctrl *gomock.Controller) *MockStream {
	mock := &MockStream{ctrl: ctrl}
	mock.recorder = &MockStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStream) EXPECT() *MockStreamMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStream) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStreamMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStream)(nil).Close))
}

// CloseRead mocks base method.
func (m *MockStream) CloseRead() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseRead")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseRead indicates an expected call of CloseRead.
func (mr *MockStreamMockRecorder) CloseRead() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseRead", reflect.TypeOf((*MockStream)(nil).CloseRead))
}

// CloseWrite mocks base method.
func (m *MockStream) CloseWrite() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseWrite")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseWrite indicates an expected call of CloseWrite.
func (mr *MockStreamMockRecorder) CloseWrite() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseWrite", reflect.TypeOf((*MockStream)(nil).CloseWrite))
}

// Conn mocks base method.
func (m *MockStream) Conn() network.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Conn")
	ret0, _ := ret[0].(network.Conn)
	return ret0
}

// Conn indicates an expected call of Conn.
func (mr *MockStreamMockRecorder) Conn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Conn", reflect.TypeOf((*MockStream)(nil).Conn))
}

// ID mocks base method.
func (m *MockStream) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockStreamMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockStream)(nil).ID))
}

// Protocol mocks base method.
func (m *MockStream) Protocol() protocol.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Protocol")
	ret0, _ := ret[0].(protocol.ID)
	return ret0
}

// Protocol indicates an expected call of Protocol.
func (mr *MockStreamMockRecorder) Protocol() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Protocol", reflect.TypeOf((*MockStream)(nil).Protocol))
}

// Read mocks base method.
func (m *MockStream) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockStreamMockRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockStream)(nil).Read), p)
}

// Reset mocks base method.
func (m *MockStream) Reset() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockStreamMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockStream)(nil).Reset))
}

// SetDeadline mocks base method.
func (m *MockStream) SetDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDeadline indicates an expected call of SetDeadline.
func (mr *MockStreamMockRecorder) SetDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*MockStream)(nil).SetDeadline), arg0)
}

// SetProtocol mocks base method.
func (m *MockStream) SetProtocol(id protocol.ID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetProtocol", id)
}

// SetProtocol indicates an expected call of SetProtocol.
func (mr *MockStreamMockRecorder) SetProtocol(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProtocol", reflect.TypeOf((*MockStream)(nil).SetProtocol), id)
}

// SetReadDeadline mocks base method.
func (m *MockStream) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline.
func (mr *MockStreamMockRecorder) SetReadDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockStream)(nil).SetReadDeadline), arg0)
}

// SetWriteDeadline mocks base method.
func (m *MockStream) SetWriteDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline.
func (mr *MockStreamMockRecorder) SetWriteDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*MockStream)(nil).SetWriteDeadline), arg0)
}

// Stat mocks base method.
func (m *MockStream) Stat() network.Stat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(network.Stat)
	return ret0
}

// Stat indicates an expected call of Stat.
func (mr *MockStreamMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockStream)(nil).Stat))
}

// Write mocks base method.
func (m *MockStream) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockStreamMockRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockStream)(nil).Write), p)
}

// MockNetwork is a mock of Network interface.
type MockNetwork struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkMockRecorder
}

// MockNetworkMockRecorder is the mock recorder for MockNetwork.
type MockNetworkMockRecorder struct {
	mock *MockNetwork
}

// NewMockNetwork creates a new mock instance.
func NewMockNetwork(ctrl *gomock.Controller) *MockNetwork {
	mock := &MockNetwork{ctrl: ctrl}
	mock.recorder = &MockNetworkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetwork) EXPECT() *MockNetworkMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockNetwork) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockNetworkMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockNetwork)(nil).Close))
}

// ClosePeer mocks base method.
func (m *MockNetwork) ClosePeer(arg0 peer.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClosePeer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClosePeer indicates an expected call of ClosePeer.
func (mr *MockNetworkMockRecorder) ClosePeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClosePeer", reflect.TypeOf((*MockNetwork)(nil).ClosePeer), arg0)
}

// Connectedness mocks base method.
func (m *MockNetwork) Connectedness(arg0 peer.ID) network.Connectedness {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connectedness", arg0)
	ret0, _ := ret[0].(network.Connectedness)
	return ret0
}

// Connectedness indicates an expected call of Connectedness.
func (mr *MockNetworkMockRecorder) Connectedness(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connectedness", reflect.TypeOf((*MockNetwork)(nil).Connectedness), arg0)
}

// Conns mocks base method.
func (m *MockNetwork) Conns() []network.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Conns")
	ret0, _ := ret[0].([]network.Conn)
	return ret0
}

// Conns indicates an expected call of Conns.
func (mr *MockNetworkMockRecorder) Conns() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Conns", reflect.TypeOf((*MockNetwork)(nil).Conns))
}

// ConnsToPeer mocks base method.
func (m *MockNetwork) ConnsToPeer(p peer.ID) []network.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnsToPeer", p)
	ret0, _ := ret[0].([]network.Conn)
	return ret0
}

// ConnsToPeer indicates an expected call of ConnsToPeer.
func (mr *MockNetworkMockRecorder) ConnsToPeer(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnsToPeer", reflect.TypeOf((*MockNetwork)(nil).ConnsToPeer), p)
}

// DialPeer mocks base method.
func (m *MockNetwork) DialPeer(arg0 context.Context, arg1 peer.ID) (network.Conn, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DialPeer", arg0, arg1)
	ret0, _ := ret[0].(network.Conn)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DialPeer indicates an expected call of DialPeer.
func (mr *MockNetworkMockRecorder) DialPeer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DialPeer", reflect.TypeOf((*MockNetwork)(nil).DialPeer), arg0, arg1)
}

// InterfaceListenAddresses mocks base method.
func (m *MockNetwork) InterfaceListenAddresses() ([]multiaddr.Multiaddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InterfaceListenAddresses")
	ret0, _ := ret[0].([]multiaddr.Multiaddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InterfaceListenAddresses indicates an expected call of InterfaceListenAddresses.
func (mr *MockNetworkMockRecorder) InterfaceListenAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InterfaceListenAddresses", reflect.TypeOf((*MockNetwork)(nil).InterfaceListenAddresses))
}

// Listen mocks base method.
func (m *MockNetwork) Listen(arg0 ...multiaddr.Multiaddr) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Listen", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Listen indicates an expected call of Listen.
func (mr *MockNetworkMockRecorder) Listen(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listen", reflect.TypeOf((*MockNetwork)(nil).Listen), arg0...)
}

// ListenAddresses mocks base method.
func (m *MockNetwork) ListenAddresses() []multiaddr.Multiaddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAddresses")
	ret0, _ := ret[0].([]multiaddr.Multiaddr)
	return ret0
}

// ListenAddresses indicates an expected call of ListenAddresses.
func (mr *MockNetworkMockRecorder) ListenAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAddresses", reflect.TypeOf((*MockNetwork)(nil).ListenAddresses))
}

// LocalPeer mocks base method.
func (m *MockNetwork) LocalPeer() peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalPeer")
	ret0, _ := ret[0].(peer.ID)
	return ret0
}

// LocalPeer indicates an expected call of LocalPeer.
func (mr *MockNetworkMockRecorder) LocalPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalPeer", reflect.TypeOf((*MockNetwork)(nil).LocalPeer))
}

// NewStream mocks base method.
func (m *MockNetwork) NewStream(arg0 context.Context, arg1 peer.ID) (network.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewStream", arg0, arg1)
	ret0, _ := ret[0].(network.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream.
func (mr *MockNetworkMockRecorder) NewStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*MockNetwork)(nil).NewStream), arg0, arg1)
}

// Notify mocks base method.
func (m *MockNetwork) Notify(arg0 network.Notifiee) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Notify", arg0)
}

// Notify indicates an expected call of Notify.
func (mr *MockNetworkMockRecorder) Notify(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockNetwork)(nil).Notify), arg0)
}

// Peers mocks base method.
func (m *MockNetwork) Peers() []peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peers")
	ret0, _ := ret[0].([]peer.ID)
	return ret0
}

// Peers indicates an expected call of Peers.
func (mr *MockNetworkMockRecorder) Peers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peers", reflect.TypeOf((*MockNetwork)(nil).Peers))
}

// Peerstore mocks base method.
func (m *MockNetwork) Peerstore() peerstore.Peerstore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peerstore")
	ret0, _ := ret[0].(peerstore.Peerstore)
	return ret0
}

// Peerstore indicates an expected call of Peerstore.
func (mr *MockNetworkMockRecorder) Peerstore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peerstore", reflect.TypeOf((*MockNetwork)(nil).Peerstore))
}

// Process mocks base method.
func (m *MockNetwork) Process() goprocess.Process {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process")
	ret0, _ := ret[0].(goprocess.Process)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockNetworkMockRecorder) Process() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockNetwork)(nil).Process))
}

// SetConnHandler mocks base method.
func (m *MockNetwork) SetConnHandler(arg0 network.ConnHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConnHandler", arg0)
}

// SetConnHandler indicates an expected call of SetConnHandler.
func (mr *MockNetworkMockRecorder) SetConnHandler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnHandler", reflect.TypeOf((*MockNetwork)(nil).SetConnHandler), arg0)
}

// SetStreamHandler mocks base method.
func (m *MockNetwork) SetStreamHandler(arg0 network.StreamHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStreamHandler", arg0)
}

// SetStreamHandler indicates an expected call of SetStreamHandler.
func (mr *MockNetworkMockRecorder) SetStreamHandler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStreamHandler", reflect.TypeOf((*MockNetwork)(nil).SetStreamHandler), arg0)
}

// StopNotify mocks base method.
func (m *MockNetwork) StopNotify(arg0 network.Notifiee) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopNotify", arg0)
}

// StopNotify indicates an expected call of StopNotify.
func (mr *MockNetworkMockRecorder) StopNotify(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopNotify", reflect.TypeOf((*MockNetwork)(nil).StopNotify), arg0)
}
