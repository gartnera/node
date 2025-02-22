// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file zetachain/zetacore/crosschain/outbound_tracker.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message zetachain.zetacore.crosschain.TxHashList
 */
export declare class TxHashList extends Message<TxHashList> {
  /**
   * @generated from field: string tx_hash = 1;
   */
  txHash: string;

  /**
   * @generated from field: string tx_signer = 2;
   */
  txSigner: string;

  /**
   * @generated from field: bool proved = 3;
   */
  proved: boolean;

  constructor(data?: PartialMessage<TxHashList>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.TxHashList";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): TxHashList;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): TxHashList;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): TxHashList;

  static equals(a: TxHashList | PlainMessage<TxHashList> | undefined, b: TxHashList | PlainMessage<TxHashList> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.crosschain.OutboundTracker
 */
export declare class OutboundTracker extends Message<OutboundTracker> {
  /**
   * format: "chain-nonce"
   *
   * @generated from field: string index = 1;
   */
  index: string;

  /**
   * @generated from field: int64 chain_id = 2;
   */
  chainId: bigint;

  /**
   * @generated from field: uint64 nonce = 3;
   */
  nonce: bigint;

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.TxHashList hash_list = 4;
   */
  hashList: TxHashList[];

  constructor(data?: PartialMessage<OutboundTracker>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.OutboundTracker";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): OutboundTracker;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): OutboundTracker;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): OutboundTracker;

  static equals(a: OutboundTracker | PlainMessage<OutboundTracker> | undefined, b: OutboundTracker | PlainMessage<OutboundTracker> | undefined): boolean;
}

