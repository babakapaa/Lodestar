// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'notes_database.dart';

// ignore_for_file: type=lint
class $NoteRowsTable extends NoteRows with TableInfo<$NoteRowsTable, NoteRow> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $NoteRowsTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<String> id = GeneratedColumn<String>(
    'id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _ownerIdMeta = const VerificationMeta(
    'ownerId',
  );
  @override
  late final GeneratedColumn<String> ownerId = GeneratedColumn<String>(
    'owner_id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: false,
    defaultValue: const Constant(''),
  );
  static const VerificationMeta _roleMeta = const VerificationMeta('role');
  @override
  late final GeneratedColumn<String> role = GeneratedColumn<String>(
    'role',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _pathMeta = const VerificationMeta('path');
  @override
  late final GeneratedColumn<String> path = GeneratedColumn<String>(
    'path',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  @override
  late final GeneratedColumnWithTypeConverter<List<String>, String> tags =
      GeneratedColumn<String>(
        'tags',
        aliasedName,
        false,
        type: DriftSqlType.string,
        requiredDuringInsert: true,
      ).withConverter<List<String>>($NoteRowsTable.$convertertags);
  @override
  late final GeneratedColumnWithTypeConverter<List<String>, String> links =
      GeneratedColumn<String>(
        'links',
        aliasedName,
        false,
        type: DriftSqlType.string,
        requiredDuringInsert: true,
      ).withConverter<List<String>>($NoteRowsTable.$converterlinks);
  static const VerificationMeta _assignmentIdMeta = const VerificationMeta(
    'assignmentId',
  );
  @override
  late final GeneratedColumn<String> assignmentId = GeneratedColumn<String>(
    'assignment_id',
    aliasedName,
    true,
    type: DriftSqlType.string,
    requiredDuringInsert: false,
  );
  static const VerificationMeta _lessonIdMeta = const VerificationMeta(
    'lessonId',
  );
  @override
  late final GeneratedColumn<String> lessonId = GeneratedColumn<String>(
    'lesson_id',
    aliasedName,
    true,
    type: DriftSqlType.string,
    requiredDuringInsert: false,
  );
  static const VerificationMeta _seqMeta = const VerificationMeta('seq');
  @override
  late final GeneratedColumn<int> seq = GeneratedColumn<int>(
    'seq',
    aliasedName,
    false,
    type: DriftSqlType.int,
    requiredDuringInsert: false,
    defaultValue: const Constant(0),
  );
  static const VerificationMeta _unsnapshottedMeta = const VerificationMeta(
    'unsnapshotted',
  );
  @override
  late final GeneratedColumn<int> unsnapshotted = GeneratedColumn<int>(
    'unsnapshotted',
    aliasedName,
    false,
    type: DriftSqlType.int,
    requiredDuringInsert: false,
    defaultValue: const Constant(0),
  );
  static const VerificationMeta _stateMeta = const VerificationMeta('state');
  @override
  late final GeneratedColumn<Uint8List> state = GeneratedColumn<Uint8List>(
    'state',
    aliasedName,
    false,
    type: DriftSqlType.blob,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _bodyMeta = const VerificationMeta('body');
  @override
  late final GeneratedColumn<String> body = GeneratedColumn<String>(
    'body',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: false,
    defaultValue: const Constant(''),
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
  static const VerificationMeta _deletedAtMeta = const VerificationMeta(
    'deletedAt',
  );
  @override
  late final GeneratedColumn<DateTime> deletedAt = GeneratedColumn<DateTime>(
    'deleted_at',
    aliasedName,
    true,
    type: DriftSqlType.dateTime,
    requiredDuringInsert: false,
  );
  static const VerificationMeta _createdMeta = const VerificationMeta(
    'created',
  );
  @override
  late final GeneratedColumn<bool> created = GeneratedColumn<bool>(
    'created',
    aliasedName,
    false,
    type: DriftSqlType.bool,
    requiredDuringInsert: false,
    defaultConstraints: GeneratedColumn.constraintIsAlways(
      'CHECK ("created" IN (0, 1))',
    ),
    defaultValue: const Constant(false),
  );
  static const VerificationMeta _metaDirtyMeta = const VerificationMeta(
    'metaDirty',
  );
  @override
  late final GeneratedColumn<bool> metaDirty = GeneratedColumn<bool>(
    'meta_dirty',
    aliasedName,
    false,
    type: DriftSqlType.bool,
    requiredDuringInsert: false,
    defaultConstraints: GeneratedColumn.constraintIsAlways(
      'CHECK ("meta_dirty" IN (0, 1))',
    ),
    defaultValue: const Constant(false),
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
    id,
    ownerId,
    role,
    path,
    tags,
    links,
    assignmentId,
    lessonId,
    seq,
    unsnapshotted,
    state,
    body,
    changedAt,
    deletedAt,
    created,
    metaDirty,
    error,
  ];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'note_rows';
  @override
  VerificationContext validateIntegrity(
    Insertable<NoteRow> instance, {
    bool isInserting = false,
  }) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    } else if (isInserting) {
      context.missing(_idMeta);
    }
    if (data.containsKey('owner_id')) {
      context.handle(
        _ownerIdMeta,
        ownerId.isAcceptableOrUnknown(data['owner_id']!, _ownerIdMeta),
      );
    }
    if (data.containsKey('role')) {
      context.handle(
        _roleMeta,
        role.isAcceptableOrUnknown(data['role']!, _roleMeta),
      );
    } else if (isInserting) {
      context.missing(_roleMeta);
    }
    if (data.containsKey('path')) {
      context.handle(
        _pathMeta,
        path.isAcceptableOrUnknown(data['path']!, _pathMeta),
      );
    } else if (isInserting) {
      context.missing(_pathMeta);
    }
    if (data.containsKey('assignment_id')) {
      context.handle(
        _assignmentIdMeta,
        assignmentId.isAcceptableOrUnknown(
          data['assignment_id']!,
          _assignmentIdMeta,
        ),
      );
    }
    if (data.containsKey('lesson_id')) {
      context.handle(
        _lessonIdMeta,
        lessonId.isAcceptableOrUnknown(data['lesson_id']!, _lessonIdMeta),
      );
    }
    if (data.containsKey('seq')) {
      context.handle(
        _seqMeta,
        seq.isAcceptableOrUnknown(data['seq']!, _seqMeta),
      );
    }
    if (data.containsKey('unsnapshotted')) {
      context.handle(
        _unsnapshottedMeta,
        unsnapshotted.isAcceptableOrUnknown(
          data['unsnapshotted']!,
          _unsnapshottedMeta,
        ),
      );
    }
    if (data.containsKey('state')) {
      context.handle(
        _stateMeta,
        state.isAcceptableOrUnknown(data['state']!, _stateMeta),
      );
    } else if (isInserting) {
      context.missing(_stateMeta);
    }
    if (data.containsKey('body')) {
      context.handle(
        _bodyMeta,
        body.isAcceptableOrUnknown(data['body']!, _bodyMeta),
      );
    }
    if (data.containsKey('changed_at')) {
      context.handle(
        _changedAtMeta,
        changedAt.isAcceptableOrUnknown(data['changed_at']!, _changedAtMeta),
      );
    } else if (isInserting) {
      context.missing(_changedAtMeta);
    }
    if (data.containsKey('deleted_at')) {
      context.handle(
        _deletedAtMeta,
        deletedAt.isAcceptableOrUnknown(data['deleted_at']!, _deletedAtMeta),
      );
    }
    if (data.containsKey('created')) {
      context.handle(
        _createdMeta,
        created.isAcceptableOrUnknown(data['created']!, _createdMeta),
      );
    }
    if (data.containsKey('meta_dirty')) {
      context.handle(
        _metaDirtyMeta,
        metaDirty.isAcceptableOrUnknown(data['meta_dirty']!, _metaDirtyMeta),
      );
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
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  NoteRow map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return NoteRow(
      id: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}id'],
      )!,
      ownerId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}owner_id'],
      )!,
      role: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}role'],
      )!,
      path: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}path'],
      )!,
      tags: $NoteRowsTable.$convertertags.fromSql(
        attachedDatabase.typeMapping.read(
          DriftSqlType.string,
          data['${effectivePrefix}tags'],
        )!,
      ),
      links: $NoteRowsTable.$converterlinks.fromSql(
        attachedDatabase.typeMapping.read(
          DriftSqlType.string,
          data['${effectivePrefix}links'],
        )!,
      ),
      assignmentId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}assignment_id'],
      ),
      lessonId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}lesson_id'],
      ),
      seq: attachedDatabase.typeMapping.read(
        DriftSqlType.int,
        data['${effectivePrefix}seq'],
      )!,
      unsnapshotted: attachedDatabase.typeMapping.read(
        DriftSqlType.int,
        data['${effectivePrefix}unsnapshotted'],
      )!,
      state: attachedDatabase.typeMapping.read(
        DriftSqlType.blob,
        data['${effectivePrefix}state'],
      )!,
      body: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}body'],
      )!,
      changedAt: attachedDatabase.typeMapping.read(
        DriftSqlType.dateTime,
        data['${effectivePrefix}changed_at'],
      )!,
      deletedAt: attachedDatabase.typeMapping.read(
        DriftSqlType.dateTime,
        data['${effectivePrefix}deleted_at'],
      ),
      created: attachedDatabase.typeMapping.read(
        DriftSqlType.bool,
        data['${effectivePrefix}created'],
      )!,
      metaDirty: attachedDatabase.typeMapping.read(
        DriftSqlType.bool,
        data['${effectivePrefix}meta_dirty'],
      )!,
      error: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}error'],
      ),
    );
  }

  @override
  $NoteRowsTable createAlias(String alias) {
    return $NoteRowsTable(attachedDatabase, alias);
  }

  static TypeConverter<List<String>, String> $convertertags =
      const StringListConverter();
  static TypeConverter<List<String>, String> $converterlinks =
      const StringListConverter();
}

class NoteRow extends DataClass implements Insertable<NoteRow> {
  final String id;
  final String ownerId;
  final String role;
  final String path;
  final List<String> tags;
  final List<String> links;
  final String? assignmentId;
  final String? lessonId;
  final int seq;
  final int unsnapshotted;
  final Uint8List state;
  final String body;
  final DateTime changedAt;
  final DateTime? deletedAt;
  final bool created;
  final bool metaDirty;
  final String? error;
  const NoteRow({
    required this.id,
    required this.ownerId,
    required this.role,
    required this.path,
    required this.tags,
    required this.links,
    this.assignmentId,
    this.lessonId,
    required this.seq,
    required this.unsnapshotted,
    required this.state,
    required this.body,
    required this.changedAt,
    this.deletedAt,
    required this.created,
    required this.metaDirty,
    this.error,
  });
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<String>(id);
    map['owner_id'] = Variable<String>(ownerId);
    map['role'] = Variable<String>(role);
    map['path'] = Variable<String>(path);
    {
      map['tags'] = Variable<String>($NoteRowsTable.$convertertags.toSql(tags));
    }
    {
      map['links'] = Variable<String>(
        $NoteRowsTable.$converterlinks.toSql(links),
      );
    }
    if (!nullToAbsent || assignmentId != null) {
      map['assignment_id'] = Variable<String>(assignmentId);
    }
    if (!nullToAbsent || lessonId != null) {
      map['lesson_id'] = Variable<String>(lessonId);
    }
    map['seq'] = Variable<int>(seq);
    map['unsnapshotted'] = Variable<int>(unsnapshotted);
    map['state'] = Variable<Uint8List>(state);
    map['body'] = Variable<String>(body);
    map['changed_at'] = Variable<DateTime>(changedAt);
    if (!nullToAbsent || deletedAt != null) {
      map['deleted_at'] = Variable<DateTime>(deletedAt);
    }
    map['created'] = Variable<bool>(created);
    map['meta_dirty'] = Variable<bool>(metaDirty);
    if (!nullToAbsent || error != null) {
      map['error'] = Variable<String>(error);
    }
    return map;
  }

  NoteRowsCompanion toCompanion(bool nullToAbsent) {
    return NoteRowsCompanion(
      id: Value(id),
      ownerId: Value(ownerId),
      role: Value(role),
      path: Value(path),
      tags: Value(tags),
      links: Value(links),
      assignmentId: assignmentId == null && nullToAbsent
          ? const Value.absent()
          : Value(assignmentId),
      lessonId: lessonId == null && nullToAbsent
          ? const Value.absent()
          : Value(lessonId),
      seq: Value(seq),
      unsnapshotted: Value(unsnapshotted),
      state: Value(state),
      body: Value(body),
      changedAt: Value(changedAt),
      deletedAt: deletedAt == null && nullToAbsent
          ? const Value.absent()
          : Value(deletedAt),
      created: Value(created),
      metaDirty: Value(metaDirty),
      error: error == null && nullToAbsent
          ? const Value.absent()
          : Value(error),
    );
  }

  factory NoteRow.fromJson(
    Map<String, dynamic> json, {
    ValueSerializer? serializer,
  }) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return NoteRow(
      id: serializer.fromJson<String>(json['id']),
      ownerId: serializer.fromJson<String>(json['ownerId']),
      role: serializer.fromJson<String>(json['role']),
      path: serializer.fromJson<String>(json['path']),
      tags: serializer.fromJson<List<String>>(json['tags']),
      links: serializer.fromJson<List<String>>(json['links']),
      assignmentId: serializer.fromJson<String?>(json['assignmentId']),
      lessonId: serializer.fromJson<String?>(json['lessonId']),
      seq: serializer.fromJson<int>(json['seq']),
      unsnapshotted: serializer.fromJson<int>(json['unsnapshotted']),
      state: serializer.fromJson<Uint8List>(json['state']),
      body: serializer.fromJson<String>(json['body']),
      changedAt: serializer.fromJson<DateTime>(json['changedAt']),
      deletedAt: serializer.fromJson<DateTime?>(json['deletedAt']),
      created: serializer.fromJson<bool>(json['created']),
      metaDirty: serializer.fromJson<bool>(json['metaDirty']),
      error: serializer.fromJson<String?>(json['error']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<String>(id),
      'ownerId': serializer.toJson<String>(ownerId),
      'role': serializer.toJson<String>(role),
      'path': serializer.toJson<String>(path),
      'tags': serializer.toJson<List<String>>(tags),
      'links': serializer.toJson<List<String>>(links),
      'assignmentId': serializer.toJson<String?>(assignmentId),
      'lessonId': serializer.toJson<String?>(lessonId),
      'seq': serializer.toJson<int>(seq),
      'unsnapshotted': serializer.toJson<int>(unsnapshotted),
      'state': serializer.toJson<Uint8List>(state),
      'body': serializer.toJson<String>(body),
      'changedAt': serializer.toJson<DateTime>(changedAt),
      'deletedAt': serializer.toJson<DateTime?>(deletedAt),
      'created': serializer.toJson<bool>(created),
      'metaDirty': serializer.toJson<bool>(metaDirty),
      'error': serializer.toJson<String?>(error),
    };
  }

  NoteRow copyWith({
    String? id,
    String? ownerId,
    String? role,
    String? path,
    List<String>? tags,
    List<String>? links,
    Value<String?> assignmentId = const Value.absent(),
    Value<String?> lessonId = const Value.absent(),
    int? seq,
    int? unsnapshotted,
    Uint8List? state,
    String? body,
    DateTime? changedAt,
    Value<DateTime?> deletedAt = const Value.absent(),
    bool? created,
    bool? metaDirty,
    Value<String?> error = const Value.absent(),
  }) => NoteRow(
    id: id ?? this.id,
    ownerId: ownerId ?? this.ownerId,
    role: role ?? this.role,
    path: path ?? this.path,
    tags: tags ?? this.tags,
    links: links ?? this.links,
    assignmentId: assignmentId.present ? assignmentId.value : this.assignmentId,
    lessonId: lessonId.present ? lessonId.value : this.lessonId,
    seq: seq ?? this.seq,
    unsnapshotted: unsnapshotted ?? this.unsnapshotted,
    state: state ?? this.state,
    body: body ?? this.body,
    changedAt: changedAt ?? this.changedAt,
    deletedAt: deletedAt.present ? deletedAt.value : this.deletedAt,
    created: created ?? this.created,
    metaDirty: metaDirty ?? this.metaDirty,
    error: error.present ? error.value : this.error,
  );
  NoteRow copyWithCompanion(NoteRowsCompanion data) {
    return NoteRow(
      id: data.id.present ? data.id.value : this.id,
      ownerId: data.ownerId.present ? data.ownerId.value : this.ownerId,
      role: data.role.present ? data.role.value : this.role,
      path: data.path.present ? data.path.value : this.path,
      tags: data.tags.present ? data.tags.value : this.tags,
      links: data.links.present ? data.links.value : this.links,
      assignmentId: data.assignmentId.present
          ? data.assignmentId.value
          : this.assignmentId,
      lessonId: data.lessonId.present ? data.lessonId.value : this.lessonId,
      seq: data.seq.present ? data.seq.value : this.seq,
      unsnapshotted: data.unsnapshotted.present
          ? data.unsnapshotted.value
          : this.unsnapshotted,
      state: data.state.present ? data.state.value : this.state,
      body: data.body.present ? data.body.value : this.body,
      changedAt: data.changedAt.present ? data.changedAt.value : this.changedAt,
      deletedAt: data.deletedAt.present ? data.deletedAt.value : this.deletedAt,
      created: data.created.present ? data.created.value : this.created,
      metaDirty: data.metaDirty.present ? data.metaDirty.value : this.metaDirty,
      error: data.error.present ? data.error.value : this.error,
    );
  }

  @override
  String toString() {
    return (StringBuffer('NoteRow(')
          ..write('id: $id, ')
          ..write('ownerId: $ownerId, ')
          ..write('role: $role, ')
          ..write('path: $path, ')
          ..write('tags: $tags, ')
          ..write('links: $links, ')
          ..write('assignmentId: $assignmentId, ')
          ..write('lessonId: $lessonId, ')
          ..write('seq: $seq, ')
          ..write('unsnapshotted: $unsnapshotted, ')
          ..write('state: $state, ')
          ..write('body: $body, ')
          ..write('changedAt: $changedAt, ')
          ..write('deletedAt: $deletedAt, ')
          ..write('created: $created, ')
          ..write('metaDirty: $metaDirty, ')
          ..write('error: $error')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(
    id,
    ownerId,
    role,
    path,
    tags,
    links,
    assignmentId,
    lessonId,
    seq,
    unsnapshotted,
    $driftBlobEquality.hash(state),
    body,
    changedAt,
    deletedAt,
    created,
    metaDirty,
    error,
  );
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is NoteRow &&
          other.id == this.id &&
          other.ownerId == this.ownerId &&
          other.role == this.role &&
          other.path == this.path &&
          other.tags == this.tags &&
          other.links == this.links &&
          other.assignmentId == this.assignmentId &&
          other.lessonId == this.lessonId &&
          other.seq == this.seq &&
          other.unsnapshotted == this.unsnapshotted &&
          $driftBlobEquality.equals(other.state, this.state) &&
          other.body == this.body &&
          other.changedAt == this.changedAt &&
          other.deletedAt == this.deletedAt &&
          other.created == this.created &&
          other.metaDirty == this.metaDirty &&
          other.error == this.error);
}

class NoteRowsCompanion extends UpdateCompanion<NoteRow> {
  final Value<String> id;
  final Value<String> ownerId;
  final Value<String> role;
  final Value<String> path;
  final Value<List<String>> tags;
  final Value<List<String>> links;
  final Value<String?> assignmentId;
  final Value<String?> lessonId;
  final Value<int> seq;
  final Value<int> unsnapshotted;
  final Value<Uint8List> state;
  final Value<String> body;
  final Value<DateTime> changedAt;
  final Value<DateTime?> deletedAt;
  final Value<bool> created;
  final Value<bool> metaDirty;
  final Value<String?> error;
  final Value<int> rowid;
  const NoteRowsCompanion({
    this.id = const Value.absent(),
    this.ownerId = const Value.absent(),
    this.role = const Value.absent(),
    this.path = const Value.absent(),
    this.tags = const Value.absent(),
    this.links = const Value.absent(),
    this.assignmentId = const Value.absent(),
    this.lessonId = const Value.absent(),
    this.seq = const Value.absent(),
    this.unsnapshotted = const Value.absent(),
    this.state = const Value.absent(),
    this.body = const Value.absent(),
    this.changedAt = const Value.absent(),
    this.deletedAt = const Value.absent(),
    this.created = const Value.absent(),
    this.metaDirty = const Value.absent(),
    this.error = const Value.absent(),
    this.rowid = const Value.absent(),
  });
  NoteRowsCompanion.insert({
    required String id,
    this.ownerId = const Value.absent(),
    required String role,
    required String path,
    required List<String> tags,
    required List<String> links,
    this.assignmentId = const Value.absent(),
    this.lessonId = const Value.absent(),
    this.seq = const Value.absent(),
    this.unsnapshotted = const Value.absent(),
    required Uint8List state,
    this.body = const Value.absent(),
    required DateTime changedAt,
    this.deletedAt = const Value.absent(),
    this.created = const Value.absent(),
    this.metaDirty = const Value.absent(),
    this.error = const Value.absent(),
    this.rowid = const Value.absent(),
  }) : id = Value(id),
       role = Value(role),
       path = Value(path),
       tags = Value(tags),
       links = Value(links),
       state = Value(state),
       changedAt = Value(changedAt);
  static Insertable<NoteRow> custom({
    Expression<String>? id,
    Expression<String>? ownerId,
    Expression<String>? role,
    Expression<String>? path,
    Expression<String>? tags,
    Expression<String>? links,
    Expression<String>? assignmentId,
    Expression<String>? lessonId,
    Expression<int>? seq,
    Expression<int>? unsnapshotted,
    Expression<Uint8List>? state,
    Expression<String>? body,
    Expression<DateTime>? changedAt,
    Expression<DateTime>? deletedAt,
    Expression<bool>? created,
    Expression<bool>? metaDirty,
    Expression<String>? error,
    Expression<int>? rowid,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (ownerId != null) 'owner_id': ownerId,
      if (role != null) 'role': role,
      if (path != null) 'path': path,
      if (tags != null) 'tags': tags,
      if (links != null) 'links': links,
      if (assignmentId != null) 'assignment_id': assignmentId,
      if (lessonId != null) 'lesson_id': lessonId,
      if (seq != null) 'seq': seq,
      if (unsnapshotted != null) 'unsnapshotted': unsnapshotted,
      if (state != null) 'state': state,
      if (body != null) 'body': body,
      if (changedAt != null) 'changed_at': changedAt,
      if (deletedAt != null) 'deleted_at': deletedAt,
      if (created != null) 'created': created,
      if (metaDirty != null) 'meta_dirty': metaDirty,
      if (error != null) 'error': error,
      if (rowid != null) 'rowid': rowid,
    });
  }

  NoteRowsCompanion copyWith({
    Value<String>? id,
    Value<String>? ownerId,
    Value<String>? role,
    Value<String>? path,
    Value<List<String>>? tags,
    Value<List<String>>? links,
    Value<String?>? assignmentId,
    Value<String?>? lessonId,
    Value<int>? seq,
    Value<int>? unsnapshotted,
    Value<Uint8List>? state,
    Value<String>? body,
    Value<DateTime>? changedAt,
    Value<DateTime?>? deletedAt,
    Value<bool>? created,
    Value<bool>? metaDirty,
    Value<String?>? error,
    Value<int>? rowid,
  }) {
    return NoteRowsCompanion(
      id: id ?? this.id,
      ownerId: ownerId ?? this.ownerId,
      role: role ?? this.role,
      path: path ?? this.path,
      tags: tags ?? this.tags,
      links: links ?? this.links,
      assignmentId: assignmentId ?? this.assignmentId,
      lessonId: lessonId ?? this.lessonId,
      seq: seq ?? this.seq,
      unsnapshotted: unsnapshotted ?? this.unsnapshotted,
      state: state ?? this.state,
      body: body ?? this.body,
      changedAt: changedAt ?? this.changedAt,
      deletedAt: deletedAt ?? this.deletedAt,
      created: created ?? this.created,
      metaDirty: metaDirty ?? this.metaDirty,
      error: error ?? this.error,
      rowid: rowid ?? this.rowid,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<String>(id.value);
    }
    if (ownerId.present) {
      map['owner_id'] = Variable<String>(ownerId.value);
    }
    if (role.present) {
      map['role'] = Variable<String>(role.value);
    }
    if (path.present) {
      map['path'] = Variable<String>(path.value);
    }
    if (tags.present) {
      map['tags'] = Variable<String>(
        $NoteRowsTable.$convertertags.toSql(tags.value),
      );
    }
    if (links.present) {
      map['links'] = Variable<String>(
        $NoteRowsTable.$converterlinks.toSql(links.value),
      );
    }
    if (assignmentId.present) {
      map['assignment_id'] = Variable<String>(assignmentId.value);
    }
    if (lessonId.present) {
      map['lesson_id'] = Variable<String>(lessonId.value);
    }
    if (seq.present) {
      map['seq'] = Variable<int>(seq.value);
    }
    if (unsnapshotted.present) {
      map['unsnapshotted'] = Variable<int>(unsnapshotted.value);
    }
    if (state.present) {
      map['state'] = Variable<Uint8List>(state.value);
    }
    if (body.present) {
      map['body'] = Variable<String>(body.value);
    }
    if (changedAt.present) {
      map['changed_at'] = Variable<DateTime>(changedAt.value);
    }
    if (deletedAt.present) {
      map['deleted_at'] = Variable<DateTime>(deletedAt.value);
    }
    if (created.present) {
      map['created'] = Variable<bool>(created.value);
    }
    if (metaDirty.present) {
      map['meta_dirty'] = Variable<bool>(metaDirty.value);
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
    return (StringBuffer('NoteRowsCompanion(')
          ..write('id: $id, ')
          ..write('ownerId: $ownerId, ')
          ..write('role: $role, ')
          ..write('path: $path, ')
          ..write('tags: $tags, ')
          ..write('links: $links, ')
          ..write('assignmentId: $assignmentId, ')
          ..write('lessonId: $lessonId, ')
          ..write('seq: $seq, ')
          ..write('unsnapshotted: $unsnapshotted, ')
          ..write('state: $state, ')
          ..write('body: $body, ')
          ..write('changedAt: $changedAt, ')
          ..write('deletedAt: $deletedAt, ')
          ..write('created: $created, ')
          ..write('metaDirty: $metaDirty, ')
          ..write('error: $error, ')
          ..write('rowid: $rowid')
          ..write(')'))
        .toString();
  }
}

class $PendingUpdatesTable extends PendingUpdates
    with TableInfo<$PendingUpdatesTable, PendingUpdate> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $PendingUpdatesTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
    'id',
    aliasedName,
    false,
    hasAutoIncrement: true,
    type: DriftSqlType.int,
    requiredDuringInsert: false,
    defaultConstraints: GeneratedColumn.constraintIsAlways(
      'PRIMARY KEY AUTOINCREMENT',
    ),
  );
  static const VerificationMeta _noteIdMeta = const VerificationMeta('noteId');
  @override
  late final GeneratedColumn<String> noteId = GeneratedColumn<String>(
    'note_id',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _keyMeta = const VerificationMeta('key');
  @override
  late final GeneratedColumn<String> key = GeneratedColumn<String>(
    'key',
    aliasedName,
    false,
    type: DriftSqlType.string,
    requiredDuringInsert: true,
  );
  static const VerificationMeta _dataMeta = const VerificationMeta('data');
  @override
  late final GeneratedColumn<Uint8List> data = GeneratedColumn<Uint8List>(
    'data',
    aliasedName,
    false,
    type: DriftSqlType.blob,
    requiredDuringInsert: true,
  );
  @override
  List<GeneratedColumn> get $columns => [id, noteId, key, data];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'pending_updates';
  @override
  VerificationContext validateIntegrity(
    Insertable<PendingUpdate> instance, {
    bool isInserting = false,
  }) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('note_id')) {
      context.handle(
        _noteIdMeta,
        noteId.isAcceptableOrUnknown(data['note_id']!, _noteIdMeta),
      );
    } else if (isInserting) {
      context.missing(_noteIdMeta);
    }
    if (data.containsKey('key')) {
      context.handle(
        _keyMeta,
        key.isAcceptableOrUnknown(data['key']!, _keyMeta),
      );
    } else if (isInserting) {
      context.missing(_keyMeta);
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
  PendingUpdate map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return PendingUpdate(
      id: attachedDatabase.typeMapping.read(
        DriftSqlType.int,
        data['${effectivePrefix}id'],
      )!,
      noteId: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}note_id'],
      )!,
      key: attachedDatabase.typeMapping.read(
        DriftSqlType.string,
        data['${effectivePrefix}key'],
      )!,
      data: attachedDatabase.typeMapping.read(
        DriftSqlType.blob,
        data['${effectivePrefix}data'],
      )!,
    );
  }

  @override
  $PendingUpdatesTable createAlias(String alias) {
    return $PendingUpdatesTable(attachedDatabase, alias);
  }
}

class PendingUpdate extends DataClass implements Insertable<PendingUpdate> {
  final int id;
  final String noteId;
  final String key;
  final Uint8List data;
  const PendingUpdate({
    required this.id,
    required this.noteId,
    required this.key,
    required this.data,
  });
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['note_id'] = Variable<String>(noteId);
    map['key'] = Variable<String>(key);
    map['data'] = Variable<Uint8List>(data);
    return map;
  }

  PendingUpdatesCompanion toCompanion(bool nullToAbsent) {
    return PendingUpdatesCompanion(
      id: Value(id),
      noteId: Value(noteId),
      key: Value(key),
      data: Value(data),
    );
  }

  factory PendingUpdate.fromJson(
    Map<String, dynamic> json, {
    ValueSerializer? serializer,
  }) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return PendingUpdate(
      id: serializer.fromJson<int>(json['id']),
      noteId: serializer.fromJson<String>(json['noteId']),
      key: serializer.fromJson<String>(json['key']),
      data: serializer.fromJson<Uint8List>(json['data']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'noteId': serializer.toJson<String>(noteId),
      'key': serializer.toJson<String>(key),
      'data': serializer.toJson<Uint8List>(data),
    };
  }

  PendingUpdate copyWith({
    int? id,
    String? noteId,
    String? key,
    Uint8List? data,
  }) => PendingUpdate(
    id: id ?? this.id,
    noteId: noteId ?? this.noteId,
    key: key ?? this.key,
    data: data ?? this.data,
  );
  PendingUpdate copyWithCompanion(PendingUpdatesCompanion data) {
    return PendingUpdate(
      id: data.id.present ? data.id.value : this.id,
      noteId: data.noteId.present ? data.noteId.value : this.noteId,
      key: data.key.present ? data.key.value : this.key,
      data: data.data.present ? data.data.value : this.data,
    );
  }

  @override
  String toString() {
    return (StringBuffer('PendingUpdate(')
          ..write('id: $id, ')
          ..write('noteId: $noteId, ')
          ..write('key: $key, ')
          ..write('data: $data')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode =>
      Object.hash(id, noteId, key, $driftBlobEquality.hash(data));
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is PendingUpdate &&
          other.id == this.id &&
          other.noteId == this.noteId &&
          other.key == this.key &&
          $driftBlobEquality.equals(other.data, this.data));
}

class PendingUpdatesCompanion extends UpdateCompanion<PendingUpdate> {
  final Value<int> id;
  final Value<String> noteId;
  final Value<String> key;
  final Value<Uint8List> data;
  const PendingUpdatesCompanion({
    this.id = const Value.absent(),
    this.noteId = const Value.absent(),
    this.key = const Value.absent(),
    this.data = const Value.absent(),
  });
  PendingUpdatesCompanion.insert({
    this.id = const Value.absent(),
    required String noteId,
    required String key,
    required Uint8List data,
  }) : noteId = Value(noteId),
       key = Value(key),
       data = Value(data);
  static Insertable<PendingUpdate> custom({
    Expression<int>? id,
    Expression<String>? noteId,
    Expression<String>? key,
    Expression<Uint8List>? data,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (noteId != null) 'note_id': noteId,
      if (key != null) 'key': key,
      if (data != null) 'data': data,
    });
  }

  PendingUpdatesCompanion copyWith({
    Value<int>? id,
    Value<String>? noteId,
    Value<String>? key,
    Value<Uint8List>? data,
  }) {
    return PendingUpdatesCompanion(
      id: id ?? this.id,
      noteId: noteId ?? this.noteId,
      key: key ?? this.key,
      data: data ?? this.data,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (noteId.present) {
      map['note_id'] = Variable<String>(noteId.value);
    }
    if (key.present) {
      map['key'] = Variable<String>(key.value);
    }
    if (data.present) {
      map['data'] = Variable<Uint8List>(data.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('PendingUpdatesCompanion(')
          ..write('id: $id, ')
          ..write('noteId: $noteId, ')
          ..write('key: $key, ')
          ..write('data: $data')
          ..write(')'))
        .toString();
  }
}

abstract class _$NotesDatabase extends GeneratedDatabase {
  _$NotesDatabase(QueryExecutor e) : super(e);
  $NotesDatabaseManager get managers => $NotesDatabaseManager(this);
  late final $NoteRowsTable noteRows = $NoteRowsTable(this);
  late final $PendingUpdatesTable pendingUpdates = $PendingUpdatesTable(this);
  @override
  Iterable<TableInfo<Table, Object?>> get allTables =>
      allSchemaEntities.whereType<TableInfo<Table, Object?>>();
  @override
  List<DatabaseSchemaEntity> get allSchemaEntities => [
    noteRows,
    pendingUpdates,
  ];
}

typedef $$NoteRowsTableCreateCompanionBuilder = NoteRowsCompanion Function({
  required String id,
  Value<String> ownerId,
  required String role,
  required String path,
  required List<String> tags,
  required List<String> links,
  Value<String?> assignmentId,
  Value<String?> lessonId,
  Value<int> seq,
  Value<int> unsnapshotted,
  required Uint8List state,
  Value<String> body,
  required DateTime changedAt,
  Value<DateTime?> deletedAt,
  Value<bool> created,
  Value<bool> metaDirty,
  Value<String?> error,
  Value<int> rowid,
});
typedef $$NoteRowsTableUpdateCompanionBuilder = NoteRowsCompanion Function({
  Value<String> id,
  Value<String> ownerId,
  Value<String> role,
  Value<String> path,
  Value<List<String>> tags,
  Value<List<String>> links,
  Value<String?> assignmentId,
  Value<String?> lessonId,
  Value<int> seq,
  Value<int> unsnapshotted,
  Value<Uint8List> state,
  Value<String> body,
  Value<DateTime> changedAt,
  Value<DateTime?> deletedAt,
  Value<bool> created,
  Value<bool> metaDirty,
  Value<String?> error,
  Value<int> rowid,
});

class $$NoteRowsTableFilterComposer
    extends Composer<_$NotesDatabase, $NoteRowsTable> {
  $$NoteRowsTableFilterComposer({
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

  ColumnFilters<String> get ownerId => $composableBuilder(
    column: $table.ownerId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get role => $composableBuilder(
    column: $table.role,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get path => $composableBuilder(
    column: $table.path,
    builder: (column) => ColumnFilters(column),
  );

  ColumnWithTypeConverterFilters<List<String>, List<String>, String> get tags =>
      $composableBuilder(
        column: $table.tags,
        builder: (column) => ColumnWithTypeConverterFilters(column),
      );

  ColumnWithTypeConverterFilters<List<String>, List<String>, String>
  get links => $composableBuilder(
    column: $table.links,
    builder: (column) => ColumnWithTypeConverterFilters(column),
  );

  ColumnFilters<String> get assignmentId => $composableBuilder(
    column: $table.assignmentId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get lessonId => $composableBuilder(
    column: $table.lessonId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<int> get seq => $composableBuilder(
    column: $table.seq,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<int> get unsnapshotted => $composableBuilder(
    column: $table.unsnapshotted,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<Uint8List> get state => $composableBuilder(
    column: $table.state,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get body => $composableBuilder(
    column: $table.body,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<DateTime> get changedAt => $composableBuilder(
    column: $table.changedAt,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<DateTime> get deletedAt => $composableBuilder(
    column: $table.deletedAt,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<bool> get created => $composableBuilder(
    column: $table.created,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<bool> get metaDirty => $composableBuilder(
    column: $table.metaDirty,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get error => $composableBuilder(
    column: $table.error,
    builder: (column) => ColumnFilters(column),
  );
}

class $$NoteRowsTableOrderingComposer
    extends Composer<_$NotesDatabase, $NoteRowsTable> {
  $$NoteRowsTableOrderingComposer({
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

  ColumnOrderings<String> get ownerId => $composableBuilder(
    column: $table.ownerId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get role => $composableBuilder(
    column: $table.role,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get path => $composableBuilder(
    column: $table.path,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get tags => $composableBuilder(
    column: $table.tags,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get links => $composableBuilder(
    column: $table.links,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get assignmentId => $composableBuilder(
    column: $table.assignmentId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get lessonId => $composableBuilder(
    column: $table.lessonId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<int> get seq => $composableBuilder(
    column: $table.seq,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<int> get unsnapshotted => $composableBuilder(
    column: $table.unsnapshotted,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<Uint8List> get state => $composableBuilder(
    column: $table.state,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get body => $composableBuilder(
    column: $table.body,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<DateTime> get changedAt => $composableBuilder(
    column: $table.changedAt,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<DateTime> get deletedAt => $composableBuilder(
    column: $table.deletedAt,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<bool> get created => $composableBuilder(
    column: $table.created,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<bool> get metaDirty => $composableBuilder(
    column: $table.metaDirty,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get error => $composableBuilder(
    column: $table.error,
    builder: (column) => ColumnOrderings(column),
  );
}

class $$NoteRowsTableAnnotationComposer
    extends Composer<_$NotesDatabase, $NoteRowsTable> {
  $$NoteRowsTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<String> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get ownerId =>
      $composableBuilder(column: $table.ownerId, builder: (column) => column);

  GeneratedColumn<String> get role =>
      $composableBuilder(column: $table.role, builder: (column) => column);

  GeneratedColumn<String> get path =>
      $composableBuilder(column: $table.path, builder: (column) => column);

  GeneratedColumnWithTypeConverter<List<String>, String> get tags =>
      $composableBuilder(column: $table.tags, builder: (column) => column);

  GeneratedColumnWithTypeConverter<List<String>, String> get links =>
      $composableBuilder(column: $table.links, builder: (column) => column);

  GeneratedColumn<String> get assignmentId => $composableBuilder(
    column: $table.assignmentId,
    builder: (column) => column,
  );

  GeneratedColumn<String> get lessonId =>
      $composableBuilder(column: $table.lessonId, builder: (column) => column);

  GeneratedColumn<int> get seq =>
      $composableBuilder(column: $table.seq, builder: (column) => column);

  GeneratedColumn<int> get unsnapshotted => $composableBuilder(
    column: $table.unsnapshotted,
    builder: (column) => column,
  );

  GeneratedColumn<Uint8List> get state =>
      $composableBuilder(column: $table.state, builder: (column) => column);

  GeneratedColumn<String> get body =>
      $composableBuilder(column: $table.body, builder: (column) => column);

  GeneratedColumn<DateTime> get changedAt =>
      $composableBuilder(column: $table.changedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get deletedAt =>
      $composableBuilder(column: $table.deletedAt, builder: (column) => column);

  GeneratedColumn<bool> get created =>
      $composableBuilder(column: $table.created, builder: (column) => column);

  GeneratedColumn<bool> get metaDirty =>
      $composableBuilder(column: $table.metaDirty, builder: (column) => column);

  GeneratedColumn<String> get error =>
      $composableBuilder(column: $table.error, builder: (column) => column);
}

class $$NoteRowsTableTableManager
    extends
        RootTableManager<
          _$NotesDatabase,
          $NoteRowsTable,
          NoteRow,
          $$NoteRowsTableFilterComposer,
          $$NoteRowsTableOrderingComposer,
          $$NoteRowsTableAnnotationComposer,
          $$NoteRowsTableCreateCompanionBuilder,
          $$NoteRowsTableUpdateCompanionBuilder,
          (NoteRow, BaseReferences<_$NotesDatabase, $NoteRowsTable, NoteRow>),
          NoteRow,
          PrefetchHooks Function()
        > {
  $$NoteRowsTableTableManager(_$NotesDatabase db, $NoteRowsTable table)
    : super(
        TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$NoteRowsTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$NoteRowsTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$NoteRowsTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback:
              ({
                Value<String> id = const Value.absent(),
                Value<String> ownerId = const Value.absent(),
                Value<String> role = const Value.absent(),
                Value<String> path = const Value.absent(),
                Value<List<String>> tags = const Value.absent(),
                Value<List<String>> links = const Value.absent(),
                Value<String?> assignmentId = const Value.absent(),
                Value<String?> lessonId = const Value.absent(),
                Value<int> seq = const Value.absent(),
                Value<int> unsnapshotted = const Value.absent(),
                Value<Uint8List> state = const Value.absent(),
                Value<String> body = const Value.absent(),
                Value<DateTime> changedAt = const Value.absent(),
                Value<DateTime?> deletedAt = const Value.absent(),
                Value<bool> created = const Value.absent(),
                Value<bool> metaDirty = const Value.absent(),
                Value<String?> error = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => NoteRowsCompanion(
                id: id,
                ownerId: ownerId,
                role: role,
                path: path,
                tags: tags,
                links: links,
                assignmentId: assignmentId,
                lessonId: lessonId,
                seq: seq,
                unsnapshotted: unsnapshotted,
                state: state,
                body: body,
                changedAt: changedAt,
                deletedAt: deletedAt,
                created: created,
                metaDirty: metaDirty,
                error: error,
                rowid: rowid,
              ),
          createCompanionCallback:
              ({
                required String id,
                Value<String> ownerId = const Value.absent(),
                required String role,
                required String path,
                required List<String> tags,
                required List<String> links,
                Value<String?> assignmentId = const Value.absent(),
                Value<String?> lessonId = const Value.absent(),
                Value<int> seq = const Value.absent(),
                Value<int> unsnapshotted = const Value.absent(),
                required Uint8List state,
                Value<String> body = const Value.absent(),
                required DateTime changedAt,
                Value<DateTime?> deletedAt = const Value.absent(),
                Value<bool> created = const Value.absent(),
                Value<bool> metaDirty = const Value.absent(),
                Value<String?> error = const Value.absent(),
                Value<int> rowid = const Value.absent(),
              }) => NoteRowsCompanion.insert(
                id: id,
                ownerId: ownerId,
                role: role,
                path: path,
                tags: tags,
                links: links,
                assignmentId: assignmentId,
                lessonId: lessonId,
                seq: seq,
                unsnapshotted: unsnapshotted,
                state: state,
                body: body,
                changedAt: changedAt,
                deletedAt: deletedAt,
                created: created,
                metaDirty: metaDirty,
                error: error,
                rowid: rowid,
              ),
          withReferenceMapper: (p0) => p0
              .map(
                (e) => (
                  e.readTable<$NoteRowsTable, NoteRow>(table),
                  BaseReferences<_$NotesDatabase, $NoteRowsTable, NoteRow>(
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

typedef $$NoteRowsTableProcessedTableManager =
    ProcessedTableManager<
      _$NotesDatabase,
      $NoteRowsTable,
      NoteRow,
      $$NoteRowsTableFilterComposer,
      $$NoteRowsTableOrderingComposer,
      $$NoteRowsTableAnnotationComposer,
      $$NoteRowsTableCreateCompanionBuilder,
      $$NoteRowsTableUpdateCompanionBuilder,
      (NoteRow, BaseReferences<_$NotesDatabase, $NoteRowsTable, NoteRow>),
      NoteRow,
      PrefetchHooks Function()
    >;
typedef $$PendingUpdatesTableCreateCompanionBuilder =
    PendingUpdatesCompanion Function({
      Value<int> id,
      required String noteId,
      required String key,
      required Uint8List data,
    });
typedef $$PendingUpdatesTableUpdateCompanionBuilder =
    PendingUpdatesCompanion Function({
      Value<int> id,
      Value<String> noteId,
      Value<String> key,
      Value<Uint8List> data,
    });

class $$PendingUpdatesTableFilterComposer
    extends Composer<_$NotesDatabase, $PendingUpdatesTable> {
  $$PendingUpdatesTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
    column: $table.id,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get noteId => $composableBuilder(
    column: $table.noteId,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<String> get key => $composableBuilder(
    column: $table.key,
    builder: (column) => ColumnFilters(column),
  );

  ColumnFilters<Uint8List> get data => $composableBuilder(
    column: $table.data,
    builder: (column) => ColumnFilters(column),
  );
}

class $$PendingUpdatesTableOrderingComposer
    extends Composer<_$NotesDatabase, $PendingUpdatesTable> {
  $$PendingUpdatesTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
    column: $table.id,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get noteId => $composableBuilder(
    column: $table.noteId,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<String> get key => $composableBuilder(
    column: $table.key,
    builder: (column) => ColumnOrderings(column),
  );

  ColumnOrderings<Uint8List> get data => $composableBuilder(
    column: $table.data,
    builder: (column) => ColumnOrderings(column),
  );
}

class $$PendingUpdatesTableAnnotationComposer
    extends Composer<_$NotesDatabase, $PendingUpdatesTable> {
  $$PendingUpdatesTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get noteId =>
      $composableBuilder(column: $table.noteId, builder: (column) => column);

  GeneratedColumn<String> get key =>
      $composableBuilder(column: $table.key, builder: (column) => column);

  GeneratedColumn<Uint8List> get data =>
      $composableBuilder(column: $table.data, builder: (column) => column);
}

class $$PendingUpdatesTableTableManager
    extends
        RootTableManager<
          _$NotesDatabase,
          $PendingUpdatesTable,
          PendingUpdate,
          $$PendingUpdatesTableFilterComposer,
          $$PendingUpdatesTableOrderingComposer,
          $$PendingUpdatesTableAnnotationComposer,
          $$PendingUpdatesTableCreateCompanionBuilder,
          $$PendingUpdatesTableUpdateCompanionBuilder,
          (
            PendingUpdate,
            BaseReferences<
              _$NotesDatabase,
              $PendingUpdatesTable,
              PendingUpdate
            >,
          ),
          PendingUpdate,
          PrefetchHooks Function()
        > {
  $$PendingUpdatesTableTableManager(
    _$NotesDatabase db,
    $PendingUpdatesTable table,
  ) : super(
        TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$PendingUpdatesTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$PendingUpdatesTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$PendingUpdatesTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback:
              ({
                Value<int> id = const Value.absent(),
                Value<String> noteId = const Value.absent(),
                Value<String> key = const Value.absent(),
                Value<Uint8List> data = const Value.absent(),
              }) => PendingUpdatesCompanion(
                id: id,
                noteId: noteId,
                key: key,
                data: data,
              ),
          createCompanionCallback:
              ({
                Value<int> id = const Value.absent(),
                required String noteId,
                required String key,
                required Uint8List data,
              }) => PendingUpdatesCompanion.insert(
                id: id,
                noteId: noteId,
                key: key,
                data: data,
              ),
          withReferenceMapper: (p0) => p0
              .map(
                (e) => (
                  e.readTable<$PendingUpdatesTable, PendingUpdate>(table),
                  BaseReferences<
                    _$NotesDatabase,
                    $PendingUpdatesTable,
                    PendingUpdate
                  >(db, table, e),
                ),
              )
              .toList(),
          prefetchHooksCallback: null,
        ),
      );
}

typedef $$PendingUpdatesTableProcessedTableManager =
    ProcessedTableManager<
      _$NotesDatabase,
      $PendingUpdatesTable,
      PendingUpdate,
      $$PendingUpdatesTableFilterComposer,
      $$PendingUpdatesTableOrderingComposer,
      $$PendingUpdatesTableAnnotationComposer,
      $$PendingUpdatesTableCreateCompanionBuilder,
      $$PendingUpdatesTableUpdateCompanionBuilder,
      (
        PendingUpdate,
        BaseReferences<_$NotesDatabase, $PendingUpdatesTable, PendingUpdate>,
      ),
      PendingUpdate,
      PrefetchHooks Function()
    >;

class $NotesDatabaseManager {
  final _$NotesDatabase _db;
  $NotesDatabaseManager(this._db);
  $$NoteRowsTableTableManager get noteRows =>
      $$NoteRowsTableTableManager(_db, _db.noteRows);
  $$PendingUpdatesTableTableManager get pendingUpdates =>
      $$PendingUpdatesTableTableManager(_db, _db.pendingUpdates);
}
