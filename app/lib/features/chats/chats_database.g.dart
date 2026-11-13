// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'chats_database.dart';

// ignore_for_file: type=lint
class $ChatRowsTable extends ChatRows with TableInfo<$ChatRowsTable, ChatRow> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $ChatRowsTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<String> id = GeneratedColumn<String>(
    'id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _dataMeta = const VerificationMeta('data');
  @override
  late final GeneratedColumn<String> data = GeneratedColumn<String>(
    'data',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _changedAtMeta = const VerificationMeta(
    'changedAt',
  );
  @override
  late final GeneratedColumn<DateTime> changedAt = GeneratedColumn<DateTime>(
    'changed_at',
    aliasedName,
    false,
    type: DriftSqlType.dateTime,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _syncedRevMeta = const VerificationMeta(
    'syncedRev',
  );
  @override
  late final GeneratedColumn<int> syncedRev = GeneratedColumn<int>(
    'synced_rev',
    aliasedName,
    false,
    type: DriftSqlType.int,
    requiredDuringInsert: false,
    defaultValue: const Constant(-1),
  );
  static const VerificationMeta _olderMeta = const VerificationMeta('older');
  @override
  late final GeneratedColumn<bool> older = GeneratedColumn<bool>(
    'older',
    aliasedName,
    false,
    type: DriftSqlType.bool,
    requiredDuringInsert: false,
    defaultConstraints: GeneratedColumn.constraintIsAlways(
      'CHECK ("older" IN (0, 1))',
    ),
    defaultValue: const Constant(true),
  );
  @override
  List<GeneratedColumn> get $columns => [id, data, changedAt, syncedRev, older];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'chat_rows';
  @override
  VerificationContext validateIntegrity(
    Insertable<ChatRow> instance, {
    bool isInserting = false,
  }) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    } else if (isInserting) {
      context.missing(_idMeta);
    }
    if (data.containsKey('data')) {
      context.handle(
        _dataMeta,
        this.data.isAcceptableOrUnknown(data['data']!, _dataMeta),
      );
    } else if (isInserting) {
      context.missing(_dataMeta);
    }
    if (data.containsKey('changed_at')) {
      context.handle(
        _changedAtMeta,
        changedAt.isAcceptableOrUnknown(data['changed_at']!, _changedAtMeta),
      );
    } else if (isInserting) {
      context.missing(_changedAtMeta);
    }
    if (data.containsKey('synced_rev')) {
      context.handle(
        _syncedRevMeta,
        syncedRev.isAcceptableOrUnknown(data['synced_rev']!, _syncedRevMeta),
      );
    }
    if (data.containsKey('older')) {
      context.handle(
        _olderMeta,
        older.isAcceptableOrUnknown(data['older']!, _olderMeta),
      );
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  ChatRow map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return ChatRow(
      id: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}id'],
      )!,
      data: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}data'],
      )!,
      changedAt: attachedDatabase.typeMapping.read(
        DriftSqlType.dateTime,
        data['${effectivePrefix}changed_at'],
      )!,
      syncedRev: attachedDatabase.typeMapping.read(
        DriftSqlType.int,
        data['${effectivePrefix}synced_rev'],
      )!,
      older: attachedDatabase.typeMapping.read(
        DriftSqlType.bool,
        data['${effectivePrefix}older'],
      )!,
    );
  }

  @override
  $ChatRowsTable createAlias(String alias) {
    return $ChatRowsTable(attachedDatabase, alias);
  }
}

class ChatRow extends DataClass implements Insertable<ChatRow> {
  final String id;
  final String data;
  final DateTime changedAt;
  final int syncedRev;
  final bool older;
  const ChatRow({
    required this.id,
    required this.data,
    required this.changedAt,
    required this.syncedRev,
    required this.older,
  });
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<String>(id);
    map['data'] = Variable<String>(data);
    map['changed_at'] = Variable<DateTime>(changedAt);
    map['synced_rev'] = Variable<int>(syncedRev);
    map['older'] = Variable<bool>(older);
    return map;
  }

  ChatRowsCompanion toCompanion(bool nullToAbsent) {
    return ChatRowsCompanion(
      id: Value(id),
      data: Value(data),
      changedAt: Value(changedAt),
      syncedRev: Value(syncedRev),
      older: Value(older),
    );
  }

  factory ChatRow.fromJson(
    Map<String, dynamic> json, {
    ValueSerializer? serializer,
  }) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return ChatRow(
      id: serializer.fromJson<String>(json['id']),
      data: serializer.fromJson<String>(json['data']),
      changedAt: serializer.fromJson<DateTime>(json['changedAt']),
      syncedRev: serializer.fromJson<int>(json['syncedRev']),
      older: serializer.fromJson<bool>(json['older']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<String>(id),
      'data': serializer.toJson<String>(data),
      'changedAt': serializer.toJson<DateTime>(changedAt),
      'syncedRev': serializer.toJson<int>(syncedRev),
      'older': serializer.toJson<bool>(older),
    };
  }

  ChatRow copyWith({
    String? id,
    String? data,
    DateTime? changedAt,
    int? syncedRev,
    bool? older,
  }) => ChatRow(
    id: id ?? this.id,
    data: data ?? this.data,
    changedAt: changedAt ?? this.changedAt,
    syncedRev: syncedRev ?? this.syncedRev,
    older: older ?? this.older,
  );
  ChatRow copyWithCompanion(ChatRowsCompanion data) {
    return ChatRow(
      id: data.id.present ? data.id.value : this.id,
      data: data.data.present ? data.data.value : this.data,
      changedAt: data.changedAt.present ? data.changedAt.value : this.changedAt,
      syncedRev: data.syncedRev.present ? data.syncedRev.value : this.syncedRev,
      older: data.older.present ? data.older.value : this.older,
    );
  }

  @override
  String toString() {
    return (StringBuffer('ChatRow(')
          ..write('id: $id, ')
          ..write('data: $data, ')
          ..write('changedAt: $changedAt, ')
          ..write('syncedRev: $syncedRev, ')
          ..write('older: $older')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(id, data, changedAt, syncedRev, older);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is ChatRow &&
          other.id == this.id &&
          other.data == this.data &&
          other.changedAt == this.changedAt &&
          other.syncedRev == this.syncedRev &&
          other.older == this.older);
}

class ChatRowsCompanion extends UpdateCompanion<ChatRow> {
  final Value<String> id;
  final Value<String> data;
  final Value<DateTime> changedAt;
  final Value<int> syncedRev;
  final Value<bool> older;
  final Value<int> rowid;
  const ChatRowsCompanion({
    this.id = const Value.absent(),
    this.data = const Value.absent(),
    this.changedAt = const Value.absent(),
    this.syncedRev = const Value.absent(),
    this.older = const Value.absent(),
    this.rowid = const Value.absent(),
  });
  ChatRowsCompanion.insert({
    required String id,
    required String data,
    required DateTime changedAt,
    this.syncedRev = const Value.absent(),
    this.older = const Value.absent(),
    this.rowid = const Value.absent(),
  }) : id = Value(id),
       data = Value(data),
       changedAt = Value(changedAt);
  static Insertable<ChatRow> custom({
    Expression<String>? id,
    Expression<String>? data,
    Expression<DateTime>? changedAt,
    Expression<int>? syncedRev,
    Expression<bool>? older,
    Expression<int>? rowid,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (data != null) 'data': data,
      if (changedAt != null) 'changed_at': changedAt,
      if (syncedRev != null) 'synced_rev': syncedRev,
      if (older != null) 'older': older,
      if (rowid != null) 'rowid': rowid,
    });
  }

  ChatRowsCompanion copyWith({
    Value<String>? id,
    Value<String>? data,
    Value<DateTime>? changedAt,
    Value<int>? syncedRev,
    Value<bool>? older,
    Value<int>? rowid,
  }) {
    return ChatRowsCompanion(
      id: id ?? this.id,
      data: data ?? this.data,
      changedAt: changedAt ?? this.changedAt,
      syncedRev: syncedRev ?? this.syncedRev,
      older: older ?? this.older,
      rowid: rowid ?? this.rowid,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<String>(id.value);
    }
    if (data.present) {
      map['data'] = Variable<String>(data.value);
    }
    if (changedAt.present) {
      map['changed_at'] = Variable<DateTime>(changedAt.value);
    }
    if (syncedRev.present) {
      map['synced_rev'] = Variable<int>(syncedRev.value);
    }
    if (older.present) {
      map['older'] = Variable<bool>(older.value);
    }
    if (rowid.present) {
      map['rowid'] = Variable<int>(rowid.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('ChatRowsCompanion(')
          ..write('id: $id, ')
          ..write('data: $data, ')
          ..write('changedAt: $changedAt, ')
          ..write('syncedRev: $syncedRev, ')
          ..write('older: $older, ')
          ..write('rowid: $rowid')
          ..write(')'))
        .toString();
  }
}

class $MessageRowsTable extends MessageRows
    with TableInfo<$MessageRowsTable, MessageRow> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $MessageRowsTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<String> id = GeneratedColumn<String>(
    'id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _chatIdMeta = const VerificationMeta('chatId');
  @override
  late final GeneratedColumn<String> chatId = GeneratedColumn<String>(
    'chat_id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _seqMeta = const VerificationMeta('seq');
  @override
  late final GeneratedColumn<int> seq = GeneratedColumn<int>(
    'seq',
    aliasedName,
    false,
    type: DriftSqlType.int,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _dataMeta = const VerificationMeta('data');
  @override
  late final GeneratedColumn<String> data = GeneratedColumn<String>(
    'data',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  @override
  List<GeneratedColumn> get $columns => [id, chatId, seq, data];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'message_rows';
  @override
  VerificationContext validateIntegrity(
    Insertable<MessageRow> instance, {
    bool isInserting = false,
  }) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    } else if (isInserting) {
      context.missing(_idMeta);
    }
    if (data.containsKey('chat_id')) {
      context.handle(
        _chatIdMeta,
        chatId.isAcceptableOrUnknown(data['chat_id']!, _chatIdMeta),
      );
    } else if (isInserting) {
      context.missing(_chatIdMeta);
    }
    if (data.containsKey('seq')) {
      context.handle(
        _seqMeta,
        seq.isAcceptableOrUnknown(data['seq']!, _seqMeta),
      );
    } else if (isInserting) {
      context.missing(_seqMeta);
    }
    if (data.containsKey('data')) {
      context.handle(
        _dataMeta,
        this.data.isAcceptableOrUnknown(data['data']!, _dataMeta),
      );
    } else if (isInserting) {
      context.missing(_dataMeta);
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  MessageRow map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return MessageRow(
      id: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}id'],
      )!,
      chatId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}chat_id'],
      )!,
      seq: attachedDatabase.typeMapping.read(
        DriftSqlType.int,
        data['${effectivePrefix}seq'],
      )!,
      data: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}data'],
      )!,
    );
  }

  @override
  $MessageRowsTable createAlias(String alias) {
    return $MessageRowsTable(attachedDatabase, alias);
  }
}

class MessageRow extends DataClass implements Insertable<MessageRow> {
  final String id;
  final String chatId;
  final int seq;
  final String data;
  const MessageRow({
    required this.id,
    required this.chatId,
    required this.seq,
    required this.data,
  });
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<String>(id);
    map['chat_id'] = Variable<String>(chatId);
    map['seq'] = Variable<int>(seq);
    map['data'] = Variable<String>(data);
    return map;
  }

  MessageRowsCompanion toCompanion(bool nullToAbsent) {
    return MessageRowsCompanion(
      id: Value(id),
      chatId: Value(chatId),
      seq: Value(seq),
      data: Value(data),
    );
  }

  factory MessageRow.fromJson(
    Map<String, dynamic> json, {
    ValueSerializer? serializer,
  }) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return MessageRow(
      id: serializer.fromJson<String>(json['id']),
      chatId: serializer.fromJson<String>(json['chatId']),
      seq: serializer.fromJson<int>(json['seq']),
      data: serializer.fromJson<String>(json['data']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<String>(id),
      'chatId': serializer.toJson<String>(chatId),
      'seq': serializer.toJson<int>(seq),
      'data': serializer.toJson<String>(data),
    };
  }

  MessageRow copyWith({String? id, String? chatId, int? seq, String? data}) =>
      MessageRow(
        id: id ?? this.id,
        chatId: chatId ?? this.chatId,
        seq: seq ?? this.seq,
        data: data ?? this.data,
      );
  MessageRow copyWithCompanion(MessageRowsCompanion data) {
    return MessageRow(
      id: data.id.present ? data.id.value : this.id,
      chatId: data.chatId.present ? data.chatId.value : this.chatId,
      seq: data.seq.present ? data.seq.value : this.seq,
      data: data.data.present ? data.data.value : this.data,
    );
  }

  @override
  String toString() {
    return (StringBuffer('MessageRow(')
          ..write('id: $id, ')
          ..write('chatId: $chatId, ')
          ..write('seq: $seq, ')
          ..write('data: $data')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(id, chatId, seq, data);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is MessageRow &&
          other.id == this.id &&
          other.chatId == this.chatId &&
          other.seq == this.seq &&
          other.data == this.data);
}

class MessageRowsCompanion extends UpdateCompanion<MessageRow> {
  final Value<String> id;
  final Value<String> chatId;
  final Value<int> seq;
  final Value<String> data;
  final Value<int> rowid;
  const MessageRowsCompanion({
    this.id = const Value.absent(),
    this.chatId = const Value.absent(),
    this.seq = const Value.absent(),
    this.data = const Value.absent(),
    this.rowid = const Value.absent(),
  });
  MessageRowsCompanion.insert({
    required String id,
    required String chatId,
    required int seq,
    required String data,
    this.rowid = const Value.absent(),
  }) : id = Value(id),
       chatId = Value(chatId),
       seq = Value(seq),
       data = Value(data);
  static Insertable<MessageRow> custom({
    Expression<String>? id,
    Expression<String>? chatId,
    Expression<int>? seq,
    Expression<String>? data,
    Expression<int>? rowid,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (chatId != null) 'chat_id': chatId,
      if (seq != null) 'seq': seq,
      if (data != null) 'data': data,
      if (rowid != null) 'rowid': rowid,
    });
  }

  MessageRowsCompanion copyWith({
    Value<String>? id,
    Value<String>? chatId,
    Value<int>? seq,
    Value<String>? data,
    Value<int>? rowid,
  }) {
    return MessageRowsCompanion(
      id: id ?? this.id,
      chatId: chatId ?? this.chatId,
      seq: seq ?? this.seq,
      data: data ?? this.data,
      rowid: rowid ?? this.rowid,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<String>(id.value);
    }
    if (chatId.present) {
      map['chat_id'] = Variable<String>(chatId.value);
    }
    if (seq.present) {
      map['seq'] = Variable<int>(seq.value);
    }
    if (data.present) {
      map['data'] = Variable<String>(data.value);
    }
    if (rowid.present) {
      map['rowid'] = Variable<int>(rowid.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('MessageRowsCompanion(')
          ..write('id: $id, ')
          ..write('chatId: $chatId, ')
          ..write('seq: $seq, ')
          ..write('data: $data, ')
          ..write('rowid: $rowid')
          ..write(')'))
        .toString();
  }
}

class $OutboxRowsTable extends OutboxRows
    with TableInfo<$OutboxRowsTable, OutboxRow> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $OutboxRowsTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _clientIdMeta = const VerificationMeta(
    'clientId',
  );
  @override
  late final GeneratedColumn<String> clientId = GeneratedColumn<String>(
    'client_id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _chatIdMeta = const VerificationMeta('chatId');
  @override
  late final GeneratedColumn<String> chatId = GeneratedColumn<String>(
    'chat_id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _bodyMeta = const VerificationMeta('body');
  @override
  late final GeneratedColumn<String> body = GeneratedColumn<String>(
    'body',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _replyToMeta = const VerificationMeta(
    'replyTo',
  );
  @override
  late final GeneratedColumn<int> replyTo = GeneratedColumn<int>(
    'reply_to',
    aliasedName,
    true,
    type: DriftSqlType.int,
    requiredDuringInsert: false,
  );
  static const VerificationMeta _createdAtMeta = const VerificationMeta(
    'createdAt',
  );
  @override
  late final GeneratedColumn<DateTime> createdAt = GeneratedColumn<DateTime>(
    'created_at',
    aliasedName,
    false,
    type: DriftSqlType.dateTime,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _errorMeta = const VerificationMeta('error');
  @override
  late final GeneratedColumn<String> error = GeneratedColumn<String>(
    'error',
    aliasedName,
    true,
    type: DriftSqlType.string,
    requiredDuringInsert: false,
  );
  @override
  List<GeneratedColumn> get $columns => [
    clientId,
    chatId,
    body,
    replyTo,
    createdAt,
    error,
  ];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'outbox_rows';
  @override
  VerificationContext validateIntegrity(
    Insertable<OutboxRow> instance, {
    bool isInserting = false,
  }) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('client_id')) {
      context.handle(
        _clientIdMeta,
        clientId.isAcceptableOrUnknown(data['client_id']!, _clientIdMeta),
      );
    } else if (isInserting) {
      context.missing(_clientIdMeta);
    }
    if (data.containsKey('chat_id')) {
      context.handle(
        _chatIdMeta,
        chatId.isAcceptableOrUnknown(data['chat_id']!, _chatIdMeta),
      );
    } else if (isInserting) {
      context.missing(_chatIdMeta);
    }
    if (data.containsKey('body')) {
      context.handle(
        _bodyMeta,
        body.isAcceptableOrUnknown(data['body']!, _bodyMeta),
      );
    } else if (isInserting) {
      context.missing(_bodyMeta);
    }
    if (data.containsKey('reply_to')) {
      context.handle(
        _replyToMeta,
        replyTo.isAcceptableOrUnknown(data['reply_to']!, _replyToMeta),
      );
    }
    if (data.containsKey('created_at')) {
      context.handle(
        _createdAtMeta,
        createdAt.isAcceptableOrUnknown(data['created_at']!, _createdAtMeta),
      );
    } else if (isInserting) {
      context.missing(_createdAtMeta);
    }
    if (data.containsKey('error')) {
      context.handle(
        _errorMeta,
        error.isAcceptableOrUnknown(data['error']!, _errorMeta),
      );
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {clientId};
  @override
  OutboxRow map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return OutboxRow(
      clientId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}client_id'],
      )!,
      chatId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}chat_id'],
      )!,
      body: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}body'],
      )!,
      replyTo: attachedDatabase.typeMapping.read(
        DriftSqlType.int,
        data['${effectivePrefix}reply_to'],
      ),
      createdAt: attachedDatabase.typeMapping.read(
        DriftSqlType.dateTime,
        data['${effectivePrefix}created_at'],
      )!,
      error: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}error'],
      ),
    );
  }

  @override
  $OutboxRowsTable createAlias(String alias) {
    return $OutboxRowsTable(attachedDatabase, alias);
  }
}

class OutboxRow extends DataClass implements Insertable<OutboxRow> {
  final String clientId;
  final String chatId;
  final String body;
  final int? replyTo;
  final DateTime createdAt;
  final String? error;
  const OutboxRow({
    required this.clientId,
    required this.chatId,
    required this.body,
    this.replyTo,
    required this.createdAt,
    this.error,
  });
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['client_id'] = Variable<String>(clientId);
    map['chat_id'] = Variable<String>(chatId);
    map['body'] = Variable<String>(body);
    if (!nullToAbsent || replyTo != null) {
      map['reply_to'] = Variable<int>(replyTo);
    }
    map['created_at'] = Variable<DateTime>(createdAt);
    if (!nullToAbsent || error != null) {
      map['error'] = Variable<String>(error);
    }
    return map;
  }

  OutboxRowsCompanion toCompanion(bool nullToAbsent) {
    return OutboxRowsCompanion(
      clientId: Value(clientId),
      chatId: Value(chatId),
      body: Value(body),
      replyTo: replyTo == null && nullToAbsent
          ? const Value.absent()
          : Value(replyTo),
      createdAt: Value(createdAt),
      error: error == null && nullToAbsent
          ? const Value.absent()
          : Value(error),
    );
  }

  factory OutboxRow.fromJson(
    Map<String, dynamic> json, {
    ValueSerializer? serializer,
  }) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return OutboxRow(
      clientId: serializer.fromJson<String>(json['clientId']),
      chatId: serializer.fromJson<String>(json['chatId']),
      body: serializer.fromJson<String>(json['body']),
      replyTo: serializer.fromJson<int?>(json['replyTo']),
      createdAt: serializer.fromJson<DateTime>(json['createdAt']),
      error: serializer.fromJson<String?>(json['error']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'clientId': serializer.toJson<String>(clientId),
      'chatId': serializer.toJson<String>(chatId),
      'body': serializer.toJson<String>(body),
      'replyTo': serializer.toJson<int?>(replyTo),
      'createdAt': serializer.toJson<DateTime>(createdAt),
      'error': serializer.toJson<String?>(error),
    };
  }

  OutboxRow copyWith({
    String? clientId,
    String? chatId,
    String? body,
    Value<int?> replyTo = const Value.absent(),
    DateTime? createdAt,
    Value<String?> error = const Value.absent(),
  }) => OutboxRow(
    clientId: clientId ?? this.clientId,
    chatId: chatId ?? this.chatId,
    body: body ?? this.body,
    replyTo: replyTo.present ? replyTo.value : this.replyTo,
    createdAt: createdAt ?? this.createdAt,
    error: error.present ? error.value : this.error,
  );
  OutboxRow copyWithCompanion(OutboxRowsCompanion data) {
    return OutboxRow(
      clientId: data.clientId.present ? data.clientId.value : this.clientId,
      chatId: data.chatId.present ? data.chatId.value : this.chatId,
      body: data.body.present ? data.body.value : this.body,
      replyTo: data.replyTo.present ? data.replyTo.value : this.replyTo,
      createdAt: data.createdAt.present ? data.createdAt.value : this.createdAt,
      error: data.error.present ? data.error.value : this.error,
    );
  }

  @override
  String toString() {
    return (StringBuffer('OutboxRow(')
          ..write('clientId: $clientId, ')
          ..write('chatId: $chatId, ')
          ..write('body: $body, ')
          ..write('replyTo: $replyTo, ')
          ..write('createdAt: $createdAt, ')
          ..write('error: $error')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode =>
      Object.hash(clientId, chatId, body, replyTo, createdAt, error);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is OutboxRow &&
          other.clientId == this.clientId &&
          other.chatId == this.chatId &&
          other.body == this.body &&
          other.replyTo == this.replyTo &&
          other.createdAt == this.createdAt &&
          other.error == this.error);
}

class OutboxRowsCompanion extends UpdateCompanion<OutboxRow> {
  final Value<String> clientId;
  final Value<String> chatId;
  final Value<String> body;
  final Value<int?> replyTo;
  final Value<DateTime> createdAt;
  final Value<String?> error;
  final Value<int> rowid;
  const OutboxRowsCompanion({
    this.clientId = const Value.absent(),
    this.chatId = const Value.absent(),
    this.body = const Value.absent(),
    this.replyTo = const Value.absent(),
    this.createdAt = const Value.absent(),
    this.error = const Value.absent(),
    this.rowid = const Value.absent(),
  });
  OutboxRowsCompanion.insert({
    required String clientId,
    required String chatId,
    required String body,
    this.replyTo = const Value.absent(),
    required DateTime createdAt,
    this.error = const Value.absent(),
    this.rowid = const Value.absent(),
  }) : clientId = Value(clientId),
       chatId = Value(chatId),
       body = Value(body),
       createdAt = Value(createdAt);
  static Insertable<OutboxRow> custom({
    Expression<String>? clientId,
    Expression<String>? chatId,
    Expression<String>? body,
    Expression<int>? replyTo,
    Expression<DateTime>? createdAt,
    Expression<String>? error,
    Expression<int>? rowid,
  }) {
    return RawValuesInsertable({
      if (clientId != null) 'client_id': clientId,
      if (chatId != null) 'chat_id': chatId,
      if (body != null) 'body': body,
      if (replyTo != null) 'reply_to': replyTo,
      if (createdAt != null) 'created_at': createdAt,
      if (error != null) 'error': error,
      if (rowid != null) 'rowid': rowid,
    });
  }

  OutboxRowsCompanion copyWith({
    Value<String>? clientId,
    Value<String>? chatId,
    Value<String>? body,
    Value<int?>? replyTo,
    Value<DateTime>? createdAt,
    Value<String?>? error,
    Value<int>? rowid,
  }) {
    return OutboxRowsCompanion(
      clientId: clientId ?? this.clientId,
      chatId: chatId ?? this.chatId,
      body: body ?? this.body,
      replyTo: replyTo ?? this.replyTo,
      createdAt: createdAt ?? this.createdAt,
      error: error ?? this.error,
      rowid: rowid ?? this.rowid,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (clientId.present) {
      map['client_id'] = Variable<String>(clientId.value);
    }
    if (chatId.present) {
      map['chat_id'] = Variable<String>(chatId.value);
    }
    if (body.present) {
      map['body'] = Variable<String>(body.value);
    }
    if (replyTo.present) {
      map['reply_to'] = Variable<int>(replyTo.value);
    }
    if (createdAt.present) {
      map['created_at'] = Variable<DateTime>(createdAt.value);
    }
    if (error.present) {
      map['error'] = Variable<String>(error.value);
    }
    if (rowid.present) {
      map['rowid'] = Variable<int>(rowid.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('OutboxRowsCompanion(')
          ..write('clientId: $clientId, ')
          ..write('chatId: $chatId, ')
          ..write('body: $body, ')
          ..write('replyTo: $replyTo, ')
          ..write('createdAt: $createdAt, ')
          ..write('error: $error, ')
          ..write('rowid: $rowid')
          ..write(')'))
        .toString();
  }
}

class $ChatValuesTable extends ChatValues
    with TableInfo<$ChatValuesTable, ChatValue> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $ChatValuesTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _keyMeta = const VerificationMeta('key');
  @override
  late final GeneratedColumn<String> key = GeneratedColumn<String>(
    'key',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _valueMeta = const VerificationMeta('value');
  @override
  late final GeneratedColumn<String> value = GeneratedColumn<String>(
    'value',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  @override
  List<GeneratedColumn> get $columns => [key, value];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'chat_values';
  @override
  VerificationContext validateIntegrity(
    Insertable<ChatValue> instance, {
    bool isInserting = false,
  }) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('key')) {
      context.handle(
        _keyMeta,
        key.isAcceptableOrUnknown(data['key']!, _keyMeta),
      );
    } else if (isInserting) {
      context.missing(_keyMeta);
    }
    if (data.containsKey('value')) {
      context.handle(
        _valueMeta,
        value.isAcceptableOrUnknown(data['value']!, _valueMeta),
      );
    } else if (isInserting) {
      context.missing(_valueMeta);
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {key};
  @override
  ChatValue map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return ChatValue(
      key: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}key'],
      )!,
      value: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}value'],
      )!,
    );
  }

  @override
  $ChatValuesTable createAlias(String alias) {
    return $ChatValuesTable(attachedDatabase, alias);
  }
}

class ChatValue extends DataClass implements Insertable<ChatValue> {
  final String key;
  final String value;
  const ChatValue({required this.key, required this.value});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['key'] = Variable<String>(key);
    map['value'] = Variable<String>(value);
    return map;
  }

  ChatValuesCompanion toCompanion(bool nullToAbsent) {
    return ChatValuesCompanion(key: Value(key), value: Value(value));
  }

  factory ChatValue.fromJson(
    Map<String, dynamic> json, {
    ValueSerializer? serializer,
  }) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return ChatValue(
      key: serializer.fromJson<String>(json['key']),
      value: serializer.fromJson<String>(json['value']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'key': serializer.toJson<String>(key),
      'value': serializer.toJson<String>(value),
    };
  }

  ChatValue copyWith({String? key, String? value}) =>
      ChatValue(key: key ?? this.key, value: value ?? this.value);
  ChatValue copyWithCompanion(ChatValuesCompanion data) {
    return ChatValue(
      key: data.key.present ? data.key.value : this.key,
      value: data.value.present ? data.value.value : this.value,
    );
  }

  @override
  String toString() {
    return (StringBuffer('ChatValue(')
          ..write('key: $key, ')
          ..write('value: $value')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(key, value);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is ChatValue &&
          other.key == this.key &&
          other.value == this.value);
}

class ChatValuesCompanion extends UpdateCompanion<ChatValue> {
  final Value<String> key;
  final Value<String> value;
  final Value<int> rowid;
  const ChatValuesCompanion({
    this.key = const Value.absent(),
    this.value = const Value.absent(),
    this.rowid = const Value.absent(),
  });
  ChatValuesCompanion.insert({
    required String key,
    required String value,
    this.rowid = const Value.absent(),
  }) : key = Value(key),
       value = Value(value);
  static Insertable<ChatValue> custom({
    Expression<String>? key,
    Expression<String>? value,
    Expression<int>? rowid,
  }) {
    return RawValuesInsertable({
      if (key != null) 'key': key,
      if (value != null) 'value': value,
      if (rowid != null) 'rowid': rowid,
    });
  }

  ChatValuesCompanion copyWith({
    Value<String>? key,
    Value<String>? value,
    Value<int>? rowid,
  }) {
    return ChatValuesCompanion(
      key: key ?? this.key,
      value: value ?? this.value,
      rowid: rowid ?? this.rowid,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (key.present) {
      map['key'] = Variable<String>(key.value);
    }
    if (value.present) {
      map['value'] = Variable<String>(value.value);
    }
    if (rowid.present) {
      map['rowid'] = Variable<int>(rowid.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('ChatValuesCompanion(')
          ..write('key: $key, ')
          ..write('value: $value, ')
          ..write('rowid: $rowid')
          ..write(')'))
        .toString();
  }
}

abstract class _$ChatsDatabase extends GeneratedDatabase {
  _$ChatsDatabase(QueryExecutor e) : super(e);
  $ChatsDatabaseManager get managers => $ChatsDatabaseManager(this);
  late final $ChatRowsTable chatRows = $ChatRowsTable(this);
  late final $MessageRowsTable messageRows = $MessageRowsTable(this);
  late final $OutboxRowsTable outboxRows = $OutboxRowsTable(this);
  late final $ChatValuesTable chatValues = $ChatValuesTable(this);
  @override
  Iterable<TableInfo<Table, Object?>> get allTables =>
      allSchemaEntities.whereType<TableInfo<Table, Object?>>();
  @override
  List<DatabaseSchemaEntity> get allSchemaEntities => [
    chatRows,
    messageRows,
    outboxRows,
    chatValues,
  ];
}

typedef $$ChatRowsTableCreateCompanionBuilder = ChatRowsCompanion Function({
  required String id,
  required String data,
  required DateTime changedAt,
  Value<int> syncedRev,
  Value<bool> older,
  Value<int> rowid,
});
typedef $$ChatRowsTableUpdateCompanionBuilder = ChatRowsCompanion Function({
  Value<String> id,
  Value<String> data,
  Value<DateTime> changedAt,
  Value<int> syncedRev,
  Value<bool> older,
  Value<int> rowid,
});

class $$ChatRowsTableFilterComposer
    extends Composer<_$ChatsDatabase, $ChatRowsTable> {
  $$ChatRowsTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<String> get id => $composableBuilder(
    column: $table.id,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get data => $composableBuilder(
    column: $table.data,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<DateTime> get changedAt => $composableBuilder(
    column: $table.changedAt,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<int> get syncedRev => $composableBuilder(
    column: $table.syncedRev,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<bool> get older => $composableBuilder(
    column: $table.older,
    builder: (column) => ColumnFilters(column),
  );
}

class $$ChatRowsTableOrderingComposer
    extends Composer<_$ChatsDatabase, $ChatRowsTable> {
  $$ChatRowsTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<String> get id => $composableBuilder(
    column: $table.id,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get data => $composableBuilder(
    column: $table.data,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<DateTime> get changedAt => $composableBuilder(
    column: $table.changedAt,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<int> get syncedRev => $composableBuilder(
    column: $table.syncedRev,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<bool> get older => $composableBuilder(
    column: $table.older,
    builder: (column) => ColumnOrderings(column),
  );
}

class $$ChatRowsTableAnnotationComposer
    extends Composer<_$ChatsDatabase, $ChatRowsTable> {
  $$ChatRowsTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<String> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get data =>
      $composableBuilder(column: $table.data, builder: (column) => column);

  GeneratedColumn<DateTime> get changedAt =>
      $composableBuilder(column: $table.changedAt, builder: (column) => column);

  GeneratedColumn<int> get syncedRev =>
      $composableBuilder(column: $table.syncedRev, builder: (column) => column);

  GeneratedColumn<bool> get older =>
      $composableBuilder(column: $table.older, builder: (column) => column);
}

class $$ChatRowsTableTableManager
    extends
        RootTableManager<
          _$ChatsDatabase,
          $ChatRowsTable,
          ChatRow,
          $$ChatRowsTableFilterComposer,
          $$ChatRowsTableOrderingComposer,
          $$ChatRowsTableAnnotationComposer,
          $$ChatRowsTableCreateCompanionBuilder,
          $$ChatRowsTableUpdateCompanionBuilder,
          (ChatRow, BaseReferences<_$ChatsDatabase, $ChatRowsTable, ChatRow>),
          ChatRow,
          PrefetchHooks Function()
        > {
  $$ChatRowsTableTableManager(_$ChatsDatabase db, $ChatRowsTable table)
    : super(
        TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$ChatRowsTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$ChatRowsTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$ChatRowsTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback:
              ({
                Value<String> id = const Value.absent(),
                Value<String> data = const Value.absent(),
                Value<DateTime> changedAt = const Value.absent(),
                Value<int> syncedRev = const Value.absent(),
                Value<bool> older = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => ChatRowsCompanion(
                id: id,
                data: data,
                changedAt: changedAt,
                syncedRev: syncedRev,
                older: older,
                rowid: rowid,
              ),
          createCompanionCallback:
              ({
                required String id,
                required String data,
                required DateTime changedAt,
                Value<int> syncedRev = const Value.absent(),
                Value<bool> older = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => ChatRowsCompanion.insert(
                id: id,
                data: data,
                changedAt: changedAt,
                syncedRev: syncedRev,
                older: older,
                rowid: rowid,
              ),
          withReferenceMapper: (p0) => p0
              .map(
                (e) => (
                  e.readTable<$ChatRowsTable, ChatRow>(table),
                  BaseReferences<_$ChatsDatabase, $ChatRowsTable, ChatRow>(
                    db,
                    table,
                    e,
                  ),
                ),
              )
              .toList(),
          prefetchHooksCallback: null,
        ),
      );
}

typedef $$ChatRowsTableProcessedTableManager =
    ProcessedTableManager<
      _$ChatsDatabase,
      $ChatRowsTable,
      ChatRow,
      $$ChatRowsTableFilterComposer,
      $$ChatRowsTableOrderingComposer,
      $$ChatRowsTableAnnotationComposer,
      $$ChatRowsTableCreateCompanionBuilder,
      $$ChatRowsTableUpdateCompanionBuilder,
      (ChatRow, BaseReferences<_$ChatsDatabase, $ChatRowsTable, ChatRow>),
      ChatRow,
      PrefetchHooks Function()
    >;
typedef $$MessageRowsTableCreateCompanionBuilder =
    MessageRowsCompanion Function({
      required String id,
      required String chatId,
      required int seq,
      required String data,
      Value<int> rowid,
    });
typedef $$MessageRowsTableUpdateCompanionBuilder =
    MessageRowsCompanion Function({
      Value<String> id,
      Value<String> chatId,
      Value<int> seq,
      Value<String> data,
      Value<int> rowid,
    });

class $$MessageRowsTableFilterComposer
    extends Composer<_$ChatsDatabase, $MessageRowsTable> {
  $$MessageRowsTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<String> get id => $composableBuilder(
    column: $table.id,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get chatId => $composableBuilder(
    column: $table.chatId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<int> get seq => $composableBuilder(
    column: $table.seq,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get data => $composableBuilder(
    column: $table.data,
    builder: (column) => ColumnFilters(column),
  );
}

class $$MessageRowsTableOrderingComposer
    extends Composer<_$ChatsDatabase, $MessageRowsTable> {
  $$MessageRowsTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<String> get id => $composableBuilder(
    column: $table.id,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get chatId => $composableBuilder(
    column: $table.chatId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<int> get seq => $composableBuilder(
    column: $table.seq,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get data => $composableBuilder(
    column: $table.data,
    builder: (column) => ColumnOrderings(column),
  );
}

class $$MessageRowsTableAnnotationComposer
    extends Composer<_$ChatsDatabase, $MessageRowsTable> {
  $$MessageRowsTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<String> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get chatId =>
      $composableBuilder(column: $table.chatId, builder: (column) => column);

  GeneratedColumn<int> get seq =>
      $composableBuilder(column: $table.seq, builder: (column) => column);

  GeneratedColumn<String> get data =>
      $composableBuilder(column: $table.data, builder: (column) => column);
}

class $$MessageRowsTableTableManager
    extends
        RootTableManager<
          _$ChatsDatabase,
          $MessageRowsTable,
          MessageRow,
          $$MessageRowsTableFilterComposer,
          $$MessageRowsTableOrderingComposer,
          $$MessageRowsTableAnnotationComposer,
          $$MessageRowsTableCreateCompanionBuilder,
          $$MessageRowsTableUpdateCompanionBuilder,
          (
            MessageRow,
            BaseReferences<_$ChatsDatabase, $MessageRowsTable, MessageRow>,
          ),
          MessageRow,
          PrefetchHooks Function()
        > {
  $$MessageRowsTableTableManager(_$ChatsDatabase db, $MessageRowsTable table)
    : super(
        TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$MessageRowsTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$MessageRowsTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$MessageRowsTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback:
              ({
                Value<String> id = const Value.absent(),
                Value<String> chatId = const Value.absent(),
                Value<int> seq = const Value.absent(),
                Value<String> data = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => MessageRowsCompanion(
                id: id,
                chatId: chatId,
                seq: seq,
                data: data,
                rowid: rowid,
              ),
          createCompanionCallback:
              ({
                required String id,
                required String chatId,
                required int seq,
                required String data,
                Value<int> rowid = const Value.absent(),
              }) => MessageRowsCompanion.insert(
                id: id,
                chatId: chatId,
                seq: seq,
                data: data,
                rowid: rowid,
              ),
          withReferenceMapper: (p0) => p0
              .map(
                (e) => (
                  e.readTable<$MessageRowsTable, MessageRow>(table),
                  BaseReferences<
                    _$ChatsDatabase,
                    $MessageRowsTable,
                    MessageRow
                  >(db, table, e),
                ),
              )
              .toList(),
          prefetchHooksCallback: null,
        ),
      );
}

typedef $$MessageRowsTableProcessedTableManager =
    ProcessedTableManager<
      _$ChatsDatabase,
      $MessageRowsTable,
      MessageRow,
      $$MessageRowsTableFilterComposer,
      $$MessageRowsTableOrderingComposer,
      $$MessageRowsTableAnnotationComposer,
      $$MessageRowsTableCreateCompanionBuilder,
      $$MessageRowsTableUpdateCompanionBuilder,
      (
        MessageRow,
        BaseReferences<_$ChatsDatabase, $MessageRowsTable, MessageRow>,
      ),
      MessageRow,
      PrefetchHooks Function()
    >;
typedef $$OutboxRowsTableCreateCompanionBuilder = OutboxRowsCompanion Function({
  required String clientId,
  required String chatId,
  required String body,
  Value<int?> replyTo,
  required DateTime createdAt,
  Value<String?> error,
  Value<int> rowid,
});
typedef $$OutboxRowsTableUpdateCompanionBuilder = OutboxRowsCompanion Function({
  Value<String> clientId,
  Value<String> chatId,
  Value<String> body,
  Value<int?> replyTo,
  Value<DateTime> createdAt,
  Value<String?> error,
  Value<int> rowid,
});

class $$OutboxRowsTableFilterComposer
    extends Composer<_$ChatsDatabase, $OutboxRowsTable> {
  $$OutboxRowsTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<String> get clientId => $composableBuilder(
    column: $table.clientId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get chatId => $composableBuilder(
    column: $table.chatId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get body => $composableBuilder(
    column: $table.body,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<int> get replyTo => $composableBuilder(
    column: $table.replyTo,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<DateTime> get createdAt => $composableBuilder(
    column: $table.createdAt,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get error => $composableBuilder(
    column: $table.error,
    builder: (column) => ColumnFilters(column),
  );
}

class $$OutboxRowsTableOrderingComposer
    extends Composer<_$ChatsDatabase, $OutboxRowsTable> {
  $$OutboxRowsTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<String> get clientId => $composableBuilder(
    column: $table.clientId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get chatId => $composableBuilder(
    column: $table.chatId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get body => $composableBuilder(
    column: $table.body,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<int> get replyTo => $composableBuilder(
    column: $table.replyTo,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<DateTime> get createdAt => $composableBuilder(
    column: $table.createdAt,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get error => $composableBuilder(
    column: $table.error,
    builder: (column) => ColumnOrderings(column),
  );
}

class $$OutboxRowsTableAnnotationComposer
    extends Composer<_$ChatsDatabase, $OutboxRowsTable> {
  $$OutboxRowsTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<String> get clientId =>
      $composableBuilder(column: $table.clientId, builder: (column) => column);

  GeneratedColumn<String> get chatId =>
      $composableBuilder(column: $table.chatId, builder: (column) => column);

  GeneratedColumn<String> get body =>
      $composableBuilder(column: $table.body, builder: (column) => column);

  GeneratedColumn<int> get replyTo =>
      $composableBuilder(column: $table.replyTo, builder: (column) => column);

  GeneratedColumn<DateTime> get createdAt =>
      $composableBuilder(column: $table.createdAt, builder: (column) => column);

  GeneratedColumn<String> get error =>
      $composableBuilder(column: $table.error, builder: (column) => column);
}

class $$OutboxRowsTableTableManager
    extends
        RootTableManager<
          _$ChatsDatabase,
          $OutboxRowsTable,
          OutboxRow,
          $$OutboxRowsTableFilterComposer,
          $$OutboxRowsTableOrderingComposer,
          $$OutboxRowsTableAnnotationComposer,
          $$OutboxRowsTableCreateCompanionBuilder,
          $$OutboxRowsTableUpdateCompanionBuilder,
          (
            OutboxRow,
            BaseReferences<_$ChatsDatabase, $OutboxRowsTable, OutboxRow>,
          ),
          OutboxRow,
          PrefetchHooks Function()
        > {
  $$OutboxRowsTableTableManager(_$ChatsDatabase db, $OutboxRowsTable table)
    : super(
        TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$OutboxRowsTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$OutboxRowsTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$OutboxRowsTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback:
              ({
                Value<String> clientId = const Value.absent(),
                Value<String> chatId = const Value.absent(),
                Value<String> body = const Value.absent(),
                Value<int?> replyTo = const Value.absent(),
                Value<DateTime> createdAt = const Value.absent(),
                Value<String?> error = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => OutboxRowsCompanion(
                clientId: clientId,
                chatId: chatId,
                body: body,
                replyTo: replyTo,
                createdAt: createdAt,
                error: error,
                rowid: rowid,
              ),
          createCompanionCallback:
              ({
                required String clientId,
                required String chatId,
                required String body,
                Value<int?> replyTo = const Value.absent(),
                required DateTime createdAt,
                Value<String?> error = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => OutboxRowsCompanion.insert(
                clientId: clientId,
                chatId: chatId,
                body: body,
                replyTo: replyTo,
                createdAt: createdAt,
                error: error,
                rowid: rowid,
              ),
          withReferenceMapper: (p0) => p0
              .map(
                (e) => (
                  e.readTable<$OutboxRowsTable, OutboxRow>(table),
                  BaseReferences<_$ChatsDatabase, $OutboxRowsTable, OutboxRow>(
                    db,
                    table,
                    e,
                  ),
                ),
              )
              .toList(),
          prefetchHooksCallback: null,
        ),
      );
}

typedef $$OutboxRowsTableProcessedTableManager =
    ProcessedTableManager<
      _$ChatsDatabase,
      $OutboxRowsTable,
      OutboxRow,
      $$OutboxRowsTableFilterComposer,
      $$OutboxRowsTableOrderingComposer,
      $$OutboxRowsTableAnnotationComposer,
      $$OutboxRowsTableCreateCompanionBuilder,
      $$OutboxRowsTableUpdateCompanionBuilder,
      (OutboxRow, BaseReferences<_$ChatsDatabase, $OutboxRowsTable, OutboxRow>),
      OutboxRow,
      PrefetchHooks Function()
    >;
typedef $$ChatValuesTableCreateCompanionBuilder = ChatValuesCompanion Function({
  required String key,
  required String value,
  Value<int> rowid,
});
typedef $$ChatValuesTableUpdateCompanionBuilder = ChatValuesCompanion Function({
  Value<String> key,
  Value<String> value,
  Value<int> rowid,
});

class $$ChatValuesTableFilterComposer
    extends Composer<_$ChatsDatabase, $ChatValuesTable> {
  $$ChatValuesTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<String> get key => $composableBuilder(
    column: $table.key,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get value => $composableBuilder(
    column: $table.value,
    builder: (column) => ColumnFilters(column),
  );
}

class $$ChatValuesTableOrderingComposer
    extends Composer<_$ChatsDatabase, $ChatValuesTable> {
  $$ChatValuesTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<String> get key => $composableBuilder(
    column: $table.key,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get value => $composableBuilder(
    column: $table.value,
    builder: (column) => ColumnOrderings(column),
  );
}

class $$ChatValuesTableAnnotationComposer
    extends Composer<_$ChatsDatabase, $ChatValuesTable> {
  $$ChatValuesTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<String> get key =>
      $composableBuilder(column: $table.key, builder: (column) => column);

  GeneratedColumn<String> get value =>
      $composableBuilder(column: $table.value, builder: (column) => column);
}

class $$ChatValuesTableTableManager
    extends
        RootTableManager<
          _$ChatsDatabase,
          $ChatValuesTable,
          ChatValue,
          $$ChatValuesTableFilterComposer,
          $$ChatValuesTableOrderingComposer,
          $$ChatValuesTableAnnotationComposer,
          $$ChatValuesTableCreateCompanionBuilder,
          $$ChatValuesTableUpdateCompanionBuilder,
          (
            ChatValue,
            BaseReferences<_$ChatsDatabase, $ChatValuesTable, ChatValue>,
          ),
          ChatValue,
          PrefetchHooks Function()
        > {
  $$ChatValuesTableTableManager(_$ChatsDatabase db, $ChatValuesTable table)
    : super(
        TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$ChatValuesTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$ChatValuesTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$ChatValuesTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<String> key = const Value.absent(),
            Value<String> value = const Value.absent(),
            Value<int> rowid = const Value.absent(),
          }) => ChatValuesCompanion(key: key, value: value, rowid: rowid),
          createCompanionCallback:
              ({
                required String key,
                required String value,
                Value<int> rowid = const Value.absent(),
              }) => ChatValuesCompanion.insert(
                key: key,
                value: value,
                rowid: rowid,
              ),
          withReferenceMapper: (p0) => p0
              .map(
                (e) => (
                  e.readTable<$ChatValuesTable, ChatValue>(table),
                  BaseReferences<_$ChatsDatabase, $ChatValuesTable, ChatValue>(
                    db,
                    table,
                    e,
                  ),
                ),
              )
              .toList(),
          prefetchHooksCallback: null,
        ),
      );
}

typedef $$ChatValuesTableProcessedTableManager =
    ProcessedTableManager<
      _$ChatsDatabase,
      $ChatValuesTable,
      ChatValue,
      $$ChatValuesTableFilterComposer,
      $$ChatValuesTableOrderingComposer,
      $$ChatValuesTableAnnotationComposer,
      $$ChatValuesTableCreateCompanionBuilder,
      $$ChatValuesTableUpdateCompanionBuilder,
      (ChatValue, BaseReferences<_$ChatsDatabase, $ChatValuesTable, ChatValue>),
      ChatValue,
      PrefetchHooks Function()
    >;

class $ChatsDatabaseManager {
  final _$ChatsDatabase _db;
  $ChatsDatabaseManager(this._db);
  $$ChatRowsTableTableManager get chatRows =>
      $$ChatRowsTableTableManager(_db, _db.chatRows);
  $$MessageRowsTableTableManager get messageRows =>
      $$MessageRowsTableTableManager(_db, _db.messageRows);
  $$OutboxRowsTableTableManager get outboxRows =>
      $$OutboxRowsTableTableManager(_db, _db.outboxRows);
  $$ChatValuesTableTableManager get chatValues =>
      $$ChatValuesTableTableManager(_db, _db.chatValues);
}
