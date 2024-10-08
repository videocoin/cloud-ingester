# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: streams/private/v1/event.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from github.com.gogo.protobuf.gogoproto import gogo_pb2 as github_dot_com_dot_gogo_dot_protobuf_dot_gogoproto_dot_gogo__pb2
from github.com.videocoin.cloud_api.streams.v1 import stream_pb2 as github_dot_com_dot_videocoin_dot_cloud__api_dot_streams_dot_v1_dot_stream__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='streams/private/v1/event.proto',
  package='cloud.api.streams.private.v1',
  syntax='proto3',
  serialized_options=_b('Z\002v1\310\342\036\001\320\342\036\001\340\342\036\001\300\343\036\001\310\343\036\001'),
  serialized_pb=_b('\n\x1estreams/private/v1/event.proto\x12\x1c\x63loud.api.streams.private.v1\x1a\x1cgoogle/api/annotations.proto\x1a-github.com/gogo/protobuf/gogoproto/gogo.proto\x1a\x36github.com/videocoin/cloud-api/streams/v1/stream.proto\"\x93\x01\n\x05\x45vent\x12\x35\n\x04type\x18\x01 \x01(\x0e\x32\'.cloud.api.streams.private.v1.EventType\x12\x1f\n\tstream_id\x18\x02 \x01(\tB\x0c\xe2\xde\x1f\x08StreamID\x12\x32\n\x06status\x18\x03 \x01(\x0e\x32\".cloud.api.streams.v1.StreamStatus*\x89\x02\n\tEventType\x12,\n\x12\x45VENT_TYPE_UNKNOWN\x10\x00\x1a\x14\x8a\x9d \x10\x45ventTypeUnknown\x12*\n\x11\x45VENT_TYPE_CREATE\x10\x01\x1a\x13\x8a\x9d \x0f\x45ventTypeCreate\x12*\n\x11\x45VENT_TYPE_UPDATE\x10\x02\x1a\x13\x8a\x9d \x0f\x45ventTypeUpdate\x12*\n\x11\x45VENT_TYPE_DELETE\x10\x03\x1a\x13\x8a\x9d \x0f\x45ventTypeDelete\x12\x37\n\x18\x45VENT_TYPE_UPDATE_STATUS\x10\x04\x1a\x19\x8a\x9d \x15\x45ventTypeUpdateStatus\x1a\x11\x88\xa3\x1e\x00\xba\xa4\x1e\tEventTypeB\x18Z\x02v1\xc8\xe2\x1e\x01\xd0\xe2\x1e\x01\xe0\xe2\x1e\x01\xc0\xe3\x1e\x01\xc8\xe3\x1e\x01\x62\x06proto3')
  ,
  dependencies=[google_dot_api_dot_annotations__pb2.DESCRIPTOR,github_dot_com_dot_gogo_dot_protobuf_dot_gogoproto_dot_gogo__pb2.DESCRIPTOR,github_dot_com_dot_videocoin_dot_cloud__api_dot_streams_dot_v1_dot_stream__pb2.DESCRIPTOR,])

_EVENTTYPE = _descriptor.EnumDescriptor(
  name='EventType',
  full_name='cloud.api.streams.private.v1.EventType',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='EVENT_TYPE_UNKNOWN', index=0, number=0,
      serialized_options=_b('\212\235 \020EventTypeUnknown'),
      type=None),
    _descriptor.EnumValueDescriptor(
      name='EVENT_TYPE_CREATE', index=1, number=1,
      serialized_options=_b('\212\235 \017EventTypeCreate'),
      type=None),
    _descriptor.EnumValueDescriptor(
      name='EVENT_TYPE_UPDATE', index=2, number=2,
      serialized_options=_b('\212\235 \017EventTypeUpdate'),
      type=None),
    _descriptor.EnumValueDescriptor(
      name='EVENT_TYPE_DELETE', index=3, number=3,
      serialized_options=_b('\212\235 \017EventTypeDelete'),
      type=None),
    _descriptor.EnumValueDescriptor(
      name='EVENT_TYPE_UPDATE_STATUS', index=4, number=4,
      serialized_options=_b('\212\235 \025EventTypeUpdateStatus'),
      type=None),
  ],
  containing_type=None,
  serialized_options=_b('\210\243\036\000\272\244\036\tEventType'),
  serialized_start=348,
  serialized_end=613,
)
_sym_db.RegisterEnumDescriptor(_EVENTTYPE)

EventType = enum_type_wrapper.EnumTypeWrapper(_EVENTTYPE)
EVENT_TYPE_UNKNOWN = 0
EVENT_TYPE_CREATE = 1
EVENT_TYPE_UPDATE = 2
EVENT_TYPE_DELETE = 3
EVENT_TYPE_UPDATE_STATUS = 4



_EVENT = _descriptor.Descriptor(
  name='Event',
  full_name='cloud.api.streams.private.v1.Event',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='type', full_name='cloud.api.streams.private.v1.Event.type', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='stream_id', full_name='cloud.api.streams.private.v1.Event.stream_id', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=_b('\342\336\037\010StreamID'), file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='status', full_name='cloud.api.streams.private.v1.Event.status', index=2,
      number=3, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=198,
  serialized_end=345,
)

_EVENT.fields_by_name['type'].enum_type = _EVENTTYPE
_EVENT.fields_by_name['status'].enum_type = github_dot_com_dot_videocoin_dot_cloud__api_dot_streams_dot_v1_dot_stream__pb2._STREAMSTATUS
DESCRIPTOR.message_types_by_name['Event'] = _EVENT
DESCRIPTOR.enum_types_by_name['EventType'] = _EVENTTYPE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Event = _reflection.GeneratedProtocolMessageType('Event', (_message.Message,), {
  'DESCRIPTOR' : _EVENT,
  '__module__' : 'streams.private.v1.event_pb2'
  # @@protoc_insertion_point(class_scope:cloud.api.streams.private.v1.Event)
  })
_sym_db.RegisterMessage(Event)


DESCRIPTOR._options = None
_EVENTTYPE._options = None
_EVENTTYPE.values_by_name["EVENT_TYPE_UNKNOWN"]._options = None
_EVENTTYPE.values_by_name["EVENT_TYPE_CREATE"]._options = None
_EVENTTYPE.values_by_name["EVENT_TYPE_UPDATE"]._options = None
_EVENTTYPE.values_by_name["EVENT_TYPE_DELETE"]._options = None
_EVENTTYPE.values_by_name["EVENT_TYPE_UPDATE_STATUS"]._options = None
_EVENT.fields_by_name['stream_id']._options = None
# @@protoc_insertion_point(module_scope)
