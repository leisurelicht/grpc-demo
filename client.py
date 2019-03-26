import time

import grpc

from protobuf import auth_pb2, auth_pb2_grpc


def gen():
    while True:
        username = input("请输入用户名: ")
        password = input("请输入密码: ")

        yield auth_pb2.Request(username=username, password=password)
        time.sleep(1)


def main():
    channel = grpc.insecure_channel('127.0.0.1:12345')
    auth_stub = auth_pb2_grpc.AUTHStub(channel)
    res = auth_stub.AuthLogin(gen())
    try:
        for r in res:
            result = r.result
            print(f"Base64编码后的结果为: {result}")
            print("-" * 10)
    except grpc._channel._Rendezvous as err:
        print(err)


if __name__ == '__main__':
    main()
