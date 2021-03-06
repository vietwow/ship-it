package k8s

import (
	"context"
	"fmt"

	shipitv1beta1 "ship-it-operator/api/v1beta1"
	"ship-it/internal/image"
	"ship-it/internal/unstructured"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	toolscache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

const imageRepositoriesIndex = "ImageRepositoriesIndex"

func eventHandler(indexer toolscache.Indexer) *toolscache.ResourceEventHandlerFuncs {
	return &toolscache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			indexer.Add(obj)
		},
		DeleteFunc: func(obj interface{}) {
			indexer.Delete(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			indexer.Delete(oldObj)
			indexer.Add(newObj)
		},
	}
}

// ImageRepositoryInformer is an informer that reconciles object events with a
// cache of HelmRelease objects and indexes their image repositories. It's used
// for querying which releases are dependant on a specific image repository.
type ImageRepositoryInformer struct {
	indexer toolscache.Indexer
}

func NewInformer(ctx context.Context, ns string) (*ImageRepositoryInformer, error) {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	c, err := cache.New(config, cache.Options{
		Scheme:    scheme,
		Namespace: ns,
	})
	if err != nil {
		return nil, err
	}

	return NewInformerWithCache(ctx, c)
}

func NewInformerWithCache(ctx context.Context, c cache.Cache) (*ImageRepositoryInformer, error) {
	informer, err := c.GetInformerForKind(shipitv1beta1.Kind("HelmRelease"))
	if err != nil {
		return nil, err
	}

	indexer := toolscache.NewIndexer(
		toolscache.DeletionHandlingMetaNamespaceKeyFunc,
		toolscache.Indexers{
			imageRepositoriesIndex: imageRepositoriesIndexFunc,
		},
	)

	informer.AddEventHandler(eventHandler(indexer))

	go c.Start(ctx.Done())

	if !c.WaitForCacheSync(ctx.Done()) {
		return nil, errors.New("repository informer: image repository cache sync failed")
	}

	return &ImageRepositoryInformer{indexer}, nil
}

// Lookup returns the namespaced names of all releases that depend on the given
// image repository URI.
func (i *ImageRepositoryInformer) Lookup(image *image.Ref) ([]types.NamespacedName, error) {
	URI := image.URI()
	objs, err := i.indexer.ByIndex(imageRepositoriesIndex, URI)
	if err != nil {
		return nil, errors.Wrapf(err, "repository informer: failed to lookup releases for image repository \"%s\"", URI)
	}

	names := make([]types.NamespacedName, 0, len(objs))

	for _, obj := range objs {
		if hr, ok := obj.(*shipitv1beta1.HelmRelease); ok {
			names = append(names, types.NamespacedName{
				Name:      hr.GetName(),
				Namespace: hr.GetNamespace(),
			})
		}
	}

	return names, nil
}

func imageRepositoriesIndexFunc(obj interface{}) ([]string, error) {
	if hr, ok := obj.(*shipitv1beta1.HelmRelease); ok {
		return imageRepositories(hr), nil
	}
	return nil, fmt.Errorf("repository informer: unexpected object type %T", obj)
}

func imageRepositories(hr *shipitv1beta1.HelmRelease) []string {
	var repos []string

	unstructured.FindAll(hr.HelmValues(), "image", func(x interface{}) {
		if img, ok := x.(map[string]interface{}); ok {
			if repo, ok := img["repository"].(string); ok {
				repos = append(repos, repo)
			}
		}
	})

	return repos
}
