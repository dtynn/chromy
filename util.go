package chromy

import (
	"context"
	"encoding/json"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
)

type CallOption func(*runtime.CallFunctionOnArgs)

func nodeIDToRemoteObjectID(ctx context.Context, cli *cdp.Client, nodeID dom.NodeID) (runtime.RemoteObjectID, error) {
	reply, err := cli.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(nodeID))
	if err != nil {
		return "", err
	}

	if reply.Object.ObjectID == nil {
		return "", ErrUnableToResolveNode
	}

	return *reply.Object.ObjectID, nil
}

func callFuncOnRemoteObject(ctx context.Context, cli *cdp.Client, objectID runtime.RemoteObjectID, declaration string, arguments []interface{}, res interface{}, opt ...CallOption) (*runtime.ExceptionDetails, error) {
	callArgs := make([]runtime.CallArgument, 0, len(arguments))
	for _, one := range arguments {
		callArg := runtime.CallArgument{}

		switch v := one.(type) {
		case *runtime.RemoteObject:
			callArg.ObjectID = v.ObjectID

		case runtime.RemoteObject:
			callArg.ObjectID = v.ObjectID

		case *runtime.RemoteObjectID:
			callArg.ObjectID = v

		case runtime.RemoteObjectID:
			callArg.ObjectID = &v

		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}

			callArg.Value = json.RawMessage(b)
		}

		callArgs = append(callArgs, callArg)
	}

	arg := runtime.NewCallFunctionOnArgs(declaration).
		SetObjectID(objectID).
		SetArguments(callArgs)

	if res != nil {
		arg.SetReturnByValue(true)
	}

	reply, err := cli.Runtime.CallFunctionOn(ctx, arg)
	if err != nil {
		return nil, err
	}

	if res != nil {
		err = json.Unmarshal(reply.Result.Value, res)
		if err != nil {
			return nil, err
		}
	}

	return reply.ExceptionDetails, nil
}
