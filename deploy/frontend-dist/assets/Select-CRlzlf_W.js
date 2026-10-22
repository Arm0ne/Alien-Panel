import{$r as e,$t as t,An as n,Bi as r,Bn as i,Ci as a,Cn as o,Ei as s,Er as c,Fr as l,Gi as u,Gn as d,Hn as f,Ir as p,Jt as m,Lr as h,Mr as g,On as _,Qi as v,Sr as y,Tn as b,Ui as x,Vi as S,Vn as C,Vr as w,Xi as ee,Xn as T,Xr as te,Yn as E,Yt as D,Zr as O,_r as k,ai as A,an as j,br as ne,cr as re,da as M,ei as N,en as ie,gr as P,ia as F,ii as I,in as ae,ji as L,jr as oe,ki as R,la as z,mr as B,na as V,ni as H,nr as U,oi as W,or as G,pr as K,qi as se,qn as ce,ri as q,rn as le,si as J,ta as Y,un as ue,vr as de,xr as X,zi as Z}from"./router-CTymuOix.js";import{o as fe,t as pe}from"./focus-detector-CpEkjhq2.js";import{a as me,r as he,t as ge}from"./light-DL1LyZED.js";import{n as _e}from"./data-freshness-C9bCcTv3.js";var Q=`v-hidden`,ve=P(`[v-hidden]`,{display:`none!important`}),ye=Z({name:`Overflow`,props:{getCounter:Function,getTail:Function,updateCounter:Function,onUpdateCount:Function,onUpdateOverflow:Function},setup(e,{slots:t}){let n=z(null),r=z(null);function i(i){let{value:a}=n,{getCounter:o,getTail:s}=e,c;if(c=o===void 0?r.value:o(),!a||!c)return;c.hasAttribute(Q)&&c.removeAttribute(Q);let{children:l}=a;if(i.showAllItemsBeforeCalculate)for(let e of l)e.hasAttribute(Q)&&e.removeAttribute(Q);let u=a.offsetWidth,d=[],f=t.tail?s?.():null,p=f?f.offsetWidth:0,m=!1,h=a.children.length-+!!t.tail;for(let t=0;t<h-1;++t){if(t<0)continue;let n=l[t];if(m){n.hasAttribute(Q)||n.setAttribute(Q,``);continue}else n.hasAttribute(Q)&&n.removeAttribute(Q);let r=n.offsetWidth;if(p+=r,d[t]=r,p>u){let{updateCounter:n}=e;for(let r=t;r>=0;--r){let i=h-1-r;n===void 0?c.textContent=`${i}`:n(i);let a=c.offsetWidth;if(p-=d[r],p+a<=u||r===0){m=!0,t=r-1,f&&(t===-1?(f.style.maxWidth=`${u-a}px`,f.style.boxSizing=`border-box`):f.style.maxWidth=``);let{onUpdateCount:n}=e;n&&n(i);break}}}}let{onUpdateOverflow:g}=e;m?g!==void 0&&g(!0):(g!==void 0&&g(!1),c.setAttribute(Q,``))}let a=de();return ve.mount({id:`vueuc/overflow`,head:!0,anchorMetaName:k,ssr:a}),se(()=>i({showAllItemsBeforeCalculate:!1})),{selfRef:n,counterRef:r,sync:i}},render(){let{$slots:e}=this;return x(()=>this.sync({showAllItemsBeforeCalculate:!1})),r(`div`,{class:`v-overflow`,ref:`selfRef`},[v(e,`default`),e.counter?e.counter():r(`span`,{style:{display:`inline-block`},ref:`counterRef`}),e.tail?e.tail():null])}});function be(e,t){t&&(se(()=>{let{value:n}=e;n&&K.registerHandler(n,t)}),Y(e,(e,t)=>{t&&K.unregisterHandler(t)},{deep:!1}),u(()=>{let{value:t}=e;t&&K.unregisterHandler(t)}))}function $(e){let t=e.filter(e=>e!==void 0);if(t.length!==0)return t.length===1?t[0]:t=>{e.forEach(e=>{e&&e(t)})}}var xe=Z({name:`Checkmark`,render(){return r(`svg`,{xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 16 16`},r(`g`,{fill:`none`},r(`path`,{d:`M14.046 3.486a.75.75 0 0 1-.032 1.06l-7.93 7.474a.85.85 0 0 1-1.188-.022l-2.68-2.72a.75.75 0 1 1 1.068-1.053l2.234 2.267l7.468-7.038a.75.75 0 0 1 1.06.032z`,fill:`currentColor`})))}}),Se=Z({name:`NBaseSelectGroupHeader`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(){let{renderLabelRef:e,renderOptionRef:t,labelFieldRef:n,nodePropsRef:r}=S(g);return{labelField:n,nodeProps:r,renderLabel:e,renderOption:t}},render(){let{clsPrefix:e,renderLabel:t,renderOption:n,nodeProps:i,tmNode:{rawNode:a}}=this,o=i?.(a),s=t?t(a,!1):T(a[this.labelField],a,!1),c=r(`div`,Object.assign({},o,{class:[`${e}-base-select-group-header`,o?.class]}),s);return a.render?a.render({node:c,option:a}):n?n({node:c,option:a,selected:!1}):c}});function Ce(e,t){return r(a,{name:`fade-in-scale-up-transition`},{default:()=>e?r(o,{clsPrefix:t,class:`${t}-base-select-option__check`},{default:()=>r(xe)}):null})}var we=Z({name:`NBaseSelectOption`,props:{clsPrefix:{type:String,required:!0},tmNode:{type:Object,required:!0}},setup(e){let{valueRef:t,pendingTmNodeRef:n,multipleRef:r,valueSetRef:i,renderLabelRef:a,renderOptionRef:o,labelFieldRef:s,valueFieldRef:c,showCheckmarkRef:l,nodePropsRef:u,handleOptionClick:d,handleOptionMouseEnter:f}=S(g),p=w(()=>{let{value:t}=n;return t?e.tmNode.key===t.key:!1});function m(t){let{tmNode:n}=e;n.disabled||d(t,n)}function h(t){let{tmNode:n}=e;n.disabled||f(t,n)}function _(t){let{tmNode:n}=e,{value:r}=p;n.disabled||r||f(t,n)}return{multiple:r,isGrouped:w(()=>{let{tmNode:t}=e,{parent:n}=t;return n&&n.rawNode.type===`group`}),showCheckmark:l,nodeProps:u,isPending:p,isSelected:w(()=>{let{value:n}=t,{value:a}=r;if(n===null)return!1;let o=e.tmNode.rawNode[c.value];if(a){let{value:e}=i;return e.has(o)}else return n===o}),labelField:s,renderLabel:a,renderOption:o,handleMouseMove:_,handleMouseEnter:h,handleClick:m}},render(){let{clsPrefix:e,tmNode:{rawNode:t},isSelected:n,isPending:i,isGrouped:a,showCheckmark:o,nodeProps:s,renderOption:c,renderLabel:l,handleClick:u,handleMouseEnter:d,handleMouseMove:f}=this,p=Ce(n,e),m=l?[l(t,n),o&&p]:[T(t[this.labelField],t,n),o&&p],h=s?.(t),g=r(`div`,Object.assign({},h,{class:[`${e}-base-select-option`,t.class,h?.class,{[`${e}-base-select-option--disabled`]:t.disabled,[`${e}-base-select-option--selected`]:n,[`${e}-base-select-option--grouped`]:a,[`${e}-base-select-option--pending`]:i,[`${e}-base-select-option--show-checkmark`]:o}],style:[h?.style||``,t.style||``],onClick:$([u,h?.onClick]),onMouseenter:$([d,h?.onMouseenter]),onMousemove:$([f,h?.onMousemove])}),r(`div`,{class:`${e}-base-select-option__content`},m));return t.render?t.render({node:g,option:t,selected:n}):c?c({node:g,option:t,selected:n}):g}}),Te=q(`base-select-menu`,`
 line-height: 1.5;
 outline: none;
 z-index: 0;
 position: relative;
 border-radius: var(--n-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-color);
`,[q(`scrollbar`,`
 max-height: var(--n-height);
 `),q(`virtual-list`,`
 max-height: var(--n-height);
 `),q(`base-select-option`,`
 min-height: var(--n-option-height);
 font-size: var(--n-option-font-size);
 display: flex;
 align-items: center;
 `,[I(`content`,`
 z-index: 1;
 white-space: nowrap;
 text-overflow: ellipsis;
 overflow: hidden;
 `)]),q(`base-select-group-header`,`
 min-height: var(--n-option-height);
 font-size: .93em;
 display: flex;
 align-items: center;
 `),q(`base-select-menu-option-wrapper`,`
 position: relative;
 width: 100%;
 `),I(`loading, empty`,`
 display: flex;
 padding: 12px 32px;
 flex: 1;
 justify-content: center;
 `),I(`loading`,`
 color: var(--n-loading-color);
 font-size: var(--n-loading-size);
 `),I(`header`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-bottom: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),I(`action`,`
 padding: 8px var(--n-option-padding-left);
 font-size: var(--n-option-font-size);
 transition: 
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 border-top: 1px solid var(--n-action-divider-color);
 color: var(--n-action-text-color);
 `),q(`base-select-group-header`,`
 position: relative;
 cursor: default;
 padding: var(--n-option-padding);
 color: var(--n-group-header-text-color);
 `),q(`base-select-option`,`
 cursor: pointer;
 position: relative;
 padding: var(--n-option-padding);
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 box-sizing: border-box;
 color: var(--n-option-text-color);
 opacity: 1;
 `,[A(`show-checkmark`,`
 padding-right: calc(var(--n-option-padding-right) + 20px);
 `),H(`&::before`,`
 content: "";
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),H(`&:active`,`
 color: var(--n-option-text-color-pressed);
 `),A(`grouped`,`
 padding-left: calc(var(--n-option-padding-left) * 1.5);
 `),A(`pending`,[H(`&::before`,`
 background-color: var(--n-option-color-pending);
 `)]),A(`selected`,`
 color: var(--n-option-text-color-active);
 `,[H(`&::before`,`
 background-color: var(--n-option-color-active);
 `),A(`pending`,[H(`&::before`,`
 background-color: var(--n-option-color-active-pending);
 `)])]),A(`disabled`,`
 cursor: not-allowed;
 `,[W(`selected`,`
 color: var(--n-option-text-color-disabled);
 `),A(`selected`,`
 opacity: var(--n-option-opacity-disabled);
 `)]),I(`check`,`
 font-size: 16px;
 position: absolute;
 right: calc(var(--n-option-padding-right) - 4px);
 top: calc(50% - 7px);
 color: var(--n-option-check-color);
 transition: color .3s var(--n-bezier);
 `,[t({enterScale:`0.5`})])])]),Ee=Z({name:`InternalSelectMenu`,props:Object.assign(Object.assign({},b.props),{clsPrefix:{type:String,required:!0},scrollable:{type:Boolean,default:!0},treeMate:{type:Object,required:!0},multiple:Boolean,size:{type:String,default:`medium`},value:{type:[String,Number,Array],default:null},autoPending:Boolean,virtualScroll:{type:Boolean,default:!0},show:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},loading:Boolean,focusable:Boolean,renderLabel:Function,renderOption:Function,nodeProps:Function,showCheckmark:{type:Boolean,default:!0},onMousedown:Function,onScroll:Function,onFocus:Function,onBlur:Function,onKeyup:Function,onKeydown:Function,onTabOut:Function,onMouseenter:Function,onMouseleave:Function,onResize:Function,resetMenuOnOptionsChange:{type:Boolean,default:!0},inlineThemeDisabled:Boolean,scrollbarProps:Object,onToggle:Function}),setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n,mergedComponentPropsRef:r}=f(e),i=_(`InternalSelectMenu`,n,t),a=b(`InternalSelectMenu`,`-internal-select-menu`,Te,me,e,M(e,`clsPrefix`)),o=z(null),s=z(null),c=z(null),l=L(()=>e.treeMate.getFlattenedNodes()),d=L(()=>ae(l.value)),p=z(null);function m(){let{treeMate:t}=e,n=null,{value:r}=e;r===null?n=t.getFirstAvailableNode():(n=e.multiple?t.getNode((r||[])[(r||[]).length-1]):t.getNode(r),(!n||n.disabled)&&(n=t.getFirstAvailableNode())),V(n||null)}function h(){let{value:t}=p;t&&!e.treeMate.getNode(t.key)&&(p.value=null)}let v;Y(()=>e.show,t=>{t?v=Y(()=>e.treeMate,()=>{e.resetMenuOnOptionsChange?(e.autoPending?m():h(),x(H)):h()},{immediate:!0}):v?.()},{immediate:!0}),u(()=>{v?.()});let y=L(()=>te(a.value.self[J(`optionHeight`,e.size)])),S=L(()=>O(a.value.self[J(`padding`,e.size)])),w=L(()=>e.multiple&&Array.isArray(e.value)?new Set(e.value):new Set),T=L(()=>{let e=l.value;return e&&e.length===0}),E=L(()=>r?.value?.Select?.renderEmpty);function D(t){let{onToggle:n}=e;n&&n(t)}function k(t){let{onScroll:n}=e;n&&n(t)}function A(e){var t;(t=c.value)==null||t.sync(),k(e)}function j(){var e;(e=c.value)==null||e.sync()}function ne(){let{value:e}=p;return e||null}function re(e,t){t.disabled||V(t,!1)}function ie(e,t){t.disabled||D(t)}function P(t){var n;N(t,`action`)||(n=e.onKeyup)==null||n.call(e,t)}function F(t){var n;N(t,`action`)||(n=e.onKeydown)==null||n.call(e,t)}function I(t){var n;(n=e.onMousedown)==null||n.call(e,t),!e.focusable&&t.preventDefault()}function R(){let{value:e}=p;e&&V(e.getNext({loop:!0}),!0)}function B(){let{value:e}=p;e&&V(e.getPrev({loop:!0}),!0)}function V(e,t=!1){p.value=e,t&&H()}function H(){var t,n;let r=p.value;if(!r)return;let i=d.value(r.key);i!==null&&(e.virtualScroll?(t=s.value)==null||t.scrollTo({index:i}):(n=c.value)==null||n.scrollTo({index:i,elSize:y.value}))}function U(t){var n;o.value?.contains(t.target)&&((n=e.onFocus)==null||n.call(e,t))}function W(t){var n;o.value?.contains(t.relatedTarget)||(n=e.onBlur)==null||n.call(e,t)}ee(g,{handleOptionMouseEnter:re,handleOptionClick:ie,valueSetRef:w,pendingTmNodeRef:p,nodePropsRef:M(e,`nodeProps`),showCheckmarkRef:M(e,`showCheckmark`),multipleRef:M(e,`multiple`),valueRef:M(e,`value`),renderLabelRef:M(e,`renderLabel`),renderOptionRef:M(e,`renderOption`),labelFieldRef:M(e,`labelField`),valueFieldRef:M(e,`valueField`)}),ee(oe,o),se(()=>{let{value:e}=c;e&&e.sync()});let G=L(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{height:r,borderRadius:i,color:o,groupHeaderTextColor:s,actionDividerColor:c,optionTextColorPressed:l,optionTextColor:u,optionTextColorDisabled:d,optionTextColorActive:f,optionOpacityDisabled:p,optionCheckColor:m,actionTextColor:h,optionColorPending:g,optionColorActive:_,loadingColor:v,loadingSize:y,optionColorActivePending:b,[J(`optionFontSize`,t)]:x,[J(`optionHeight`,t)]:S,[J(`optionPadding`,t)]:C}}=a.value;return{"--n-height":r,"--n-action-divider-color":c,"--n-action-text-color":h,"--n-bezier":n,"--n-border-radius":i,"--n-color":o,"--n-option-font-size":x,"--n-group-header-text-color":s,"--n-option-check-color":m,"--n-option-color-pending":g,"--n-option-color-active":_,"--n-option-color-active-pending":b,"--n-option-height":S,"--n-option-opacity-disabled":p,"--n-option-text-color":u,"--n-option-text-color-active":f,"--n-option-text-color-disabled":d,"--n-option-text-color-pressed":l,"--n-option-padding":C,"--n-option-padding-left":O(C,`left`),"--n-option-padding-right":O(C,`right`),"--n-loading-color":v,"--n-loading-size":y}}),{inlineThemeDisabled:K}=e,ce=K?C(`internal-select-menu`,L(()=>e.size[0]),G,e):void 0,q={selfRef:o,next:R,prev:B,getPendingTmNode:ne};return be(o,e.onResize),Object.assign({mergedTheme:a,mergedClsPrefix:t,rtlEnabled:i,virtualListRef:s,scrollbarRef:c,itemSize:y,padding:S,flattenedNodes:l,empty:T,mergedRenderEmpty:E,virtualListContainer(){let{value:e}=s;return e?.listElRef},virtualListContent(){let{value:e}=s;return e?.itemsElRef},doScroll:k,handleFocusin:U,handleFocusout:W,handleKeyUp:P,handleKeyDown:F,handleMouseDown:I,handleVirtualListResize:j,handleVirtualListScroll:A,cssVars:K?void 0:G,themeClass:ce?.themeClass,onRender:ce?.onRender},q)},render(){let{$slots:e,virtualScroll:t,clsPrefix:n,mergedTheme:i,themeClass:a,onRender:o}=this;return o?.(),r(`div`,{ref:`selfRef`,tabindex:this.focusable?0:-1,class:[`${n}-base-select-menu`,`${n}-base-select-menu--${this.size}-size`,this.rtlEnabled&&`${n}-base-select-menu--rtl`,a,this.multiple&&`${n}-base-select-menu--multiple`],style:this.cssVars,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onKeyup:this.handleKeyUp,onKeydown:this.handleKeyDown,onMousedown:this.handleMouseDown,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},E(e.header,e=>e&&r(`div`,{class:`${n}-base-select-menu__header`,"data-header":!0,key:`header`},e)),this.loading?r(`div`,{class:`${n}-base-select-menu__loading`},r(ue,{clsPrefix:n,strokeWidth:20})):this.empty?r(`div`,{class:`${n}-base-select-menu__empty`,"data-empty":!0},ce(e.empty,()=>[this.mergedRenderEmpty?.call(this)||r(ie,{theme:i.peers.Empty,themeOverrides:i.peerOverrides.Empty,size:this.size})])):r(j,Object.assign({ref:`scrollbarRef`,theme:i.peers.Scrollbar,themeOverrides:i.peerOverrides.Scrollbar,scrollable:this.scrollable,container:t?this.virtualListContainer:void 0,content:t?this.virtualListContent:void 0,onScroll:t?void 0:this.doScroll},this.scrollbarProps),{default:()=>t?r(fe,{ref:`virtualListRef`,class:`${n}-virtual-list`,items:this.flattenedNodes,itemSize:this.itemSize,showScrollbar:!1,paddingTop:this.padding.top,paddingBottom:this.padding.bottom,onResize:this.handleVirtualListResize,onScroll:this.handleVirtualListScroll,itemResizable:!0},{default:({item:e})=>e.isGroup?r(Se,{key:e.key,clsPrefix:n,tmNode:e}):e.ignored?null:r(we,{clsPrefix:n,key:e.key,tmNode:e})}):r(`div`,{class:`${n}-base-select-menu-option-wrapper`,style:{paddingTop:this.padding.top,paddingBottom:this.padding.bottom}},this.flattenedNodes.map(e=>e.isGroup?r(Se,{key:e.key,clsPrefix:n,tmNode:e}):r(we,{clsPrefix:n,key:e.key,tmNode:e})))}),E(e.action,e=>e&&[r(`div`,{class:`${n}-base-select-menu__action`,"data-action":!0,key:`action`},e),r(pe,{onFocus:this.onTabOut,key:`focus-detector`})]))}}),De=H([q(`base-selection`,`
 --n-padding-single: var(--n-padding-single-top) var(--n-padding-single-right) var(--n-padding-single-bottom) var(--n-padding-single-left);
 --n-padding-multiple: var(--n-padding-multiple-top) var(--n-padding-multiple-right) var(--n-padding-multiple-bottom) var(--n-padding-multiple-left);
 position: relative;
 z-index: auto;
 box-shadow: none;
 width: 100%;
 max-width: 100%;
 display: inline-block;
 vertical-align: bottom;
 border-radius: var(--n-border-radius);
 min-height: var(--n-height);
 line-height: 1.5;
 font-size: var(--n-font-size);
 `,[q(`base-loading`,`
 color: var(--n-loading-color);
 `),q(`base-selection-tags`,`min-height: var(--n-height);`),I(`border, state-border`,`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border: var(--n-border);
 border-radius: inherit;
 transition:
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),I(`state-border`,`
 z-index: 1;
 border-color: #0000;
 `),q(`base-suffix`,`
 cursor: pointer;
 position: absolute;
 top: 50%;
 transform: translateY(-50%);
 right: 10px;
 `,[I(`arrow`,`
 font-size: var(--n-arrow-size);
 color: var(--n-arrow-color);
 transition: color .3s var(--n-bezier);
 `)]),q(`base-selection-overlay`,`
 display: flex;
 align-items: center;
 white-space: nowrap;
 pointer-events: none;
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 padding: var(--n-padding-single);
 transition: color .3s var(--n-bezier);
 `,[I(`wrapper`,`
 flex-basis: 0;
 flex-grow: 1;
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),q(`base-selection-placeholder`,`
 color: var(--n-placeholder-color);
 `,[I(`inner`,`
 max-width: 100%;
 overflow: hidden;
 `)]),q(`base-selection-tags`,`
 cursor: pointer;
 outline: none;
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 display: flex;
 padding: var(--n-padding-multiple);
 flex-wrap: wrap;
 align-items: center;
 width: 100%;
 vertical-align: bottom;
 background-color: var(--n-color);
 border-radius: inherit;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),q(`base-selection-label`,`
 height: var(--n-height);
 display: inline-flex;
 width: 100%;
 vertical-align: bottom;
 cursor: pointer;
 outline: none;
 z-index: auto;
 box-sizing: border-box;
 position: relative;
 transition:
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 border-radius: inherit;
 background-color: var(--n-color);
 align-items: center;
 `,[q(`base-selection-input`,`
 font-size: inherit;
 line-height: inherit;
 outline: none;
 cursor: pointer;
 box-sizing: border-box;
 border:none;
 width: 100%;
 padding: var(--n-padding-single);
 background-color: #0000;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 caret-color: var(--n-caret-color);
 `,[I(`content`,`
 text-overflow: ellipsis;
 overflow: hidden;
 white-space: nowrap; 
 `)]),I(`render-label`,`
 color: var(--n-text-color);
 `)]),W(`disabled`,[H(`&:hover`,[I(`state-border`,`
 box-shadow: var(--n-box-shadow-hover);
 border: var(--n-border-hover);
 `)]),A(`focus`,[I(`state-border`,`
 box-shadow: var(--n-box-shadow-focus);
 border: var(--n-border-focus);
 `)]),A(`active`,[I(`state-border`,`
 box-shadow: var(--n-box-shadow-active);
 border: var(--n-border-active);
 `),q(`base-selection-label`,`background-color: var(--n-color-active);`),q(`base-selection-tags`,`background-color: var(--n-color-active);`)])]),A(`disabled`,`cursor: not-allowed;`,[I(`arrow`,`
 color: var(--n-arrow-color-disabled);
 `),q(`base-selection-label`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `,[q(`base-selection-input`,`
 cursor: not-allowed;
 color: var(--n-text-color-disabled);
 `),I(`render-label`,`
 color: var(--n-text-color-disabled);
 `)]),q(`base-selection-tags`,`
 cursor: not-allowed;
 background-color: var(--n-color-disabled);
 `),q(`base-selection-placeholder`,`
 cursor: not-allowed;
 color: var(--n-placeholder-color-disabled);
 `)]),q(`base-selection-input-tag`,`
 height: calc(var(--n-height) - 6px);
 line-height: calc(var(--n-height) - 6px);
 outline: none;
 display: none;
 position: relative;
 margin-bottom: 3px;
 max-width: 100%;
 vertical-align: bottom;
 `,[I(`input`,`
 font-size: inherit;
 font-family: inherit;
 min-width: 1px;
 padding: 0;
 background-color: #0000;
 outline: none;
 border: none;
 max-width: 100%;
 overflow: hidden;
 width: 1em;
 line-height: inherit;
 cursor: pointer;
 color: var(--n-text-color);
 caret-color: var(--n-caret-color);
 `),I(`mirror`,`
 position: absolute;
 left: 0;
 top: 0;
 white-space: pre;
 visibility: hidden;
 user-select: none;
 -webkit-user-select: none;
 opacity: 0;
 `)]),[`warning`,`error`].map(e=>A(`${e}-status`,[I(`state-border`,`border: var(--n-border-${e});`),W(`disabled`,[H(`&:hover`,[I(`state-border`,`
 box-shadow: var(--n-box-shadow-hover-${e});
 border: var(--n-border-hover-${e});
 `)]),A(`active`,[I(`state-border`,`
 box-shadow: var(--n-box-shadow-active-${e});
 border: var(--n-border-active-${e});
 `),q(`base-selection-label`,`background-color: var(--n-color-active-${e});`),q(`base-selection-tags`,`background-color: var(--n-color-active-${e});`)]),A(`focus`,[I(`state-border`,`
 box-shadow: var(--n-box-shadow-focus-${e});
 border: var(--n-border-focus-${e});
 `)])])]))]),q(`base-selection-popover`,`
 margin-bottom: -3px;
 display: flex;
 flex-wrap: wrap;
 margin-right: -8px;
 `),q(`base-selection-tag-wrapper`,`
 max-width: 100%;
 display: inline-flex;
 padding: 0 7px 3px 0;
 `,[H(`&:last-child`,`padding-right: 0;`),q(`tag`,`
 font-size: 14px;
 max-width: 100%;
 `,[I(`content`,`
 line-height: 1.25;
 text-overflow: ellipsis;
 overflow: hidden;
 `)])])]),Oe=Z({name:`InternalSelection`,props:Object.assign(Object.assign({},b.props),{clsPrefix:{type:String,required:!0},bordered:{type:Boolean,default:void 0},active:Boolean,pattern:{type:String,default:``},placeholder:String,selectedOption:{type:Object,default:null},selectedOptions:{type:Array,default:null},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},multiple:Boolean,filterable:Boolean,clearable:Boolean,disabled:Boolean,size:{type:String,default:`medium`},loading:Boolean,autofocus:Boolean,showArrow:{type:Boolean,default:!0},inputProps:Object,focused:Boolean,renderTag:Function,onKeydown:Function,onClick:Function,onBlur:Function,onFocus:Function,onDeleteOption:Function,maxTagCount:[String,Number],ellipsisTagPopoverProps:Object,onClear:Function,onPatternInput:Function,onPatternFocus:Function,onPatternBlur:Function,renderLabel:Function,status:String,inlineThemeDisabled:Boolean,ignoreComposition:{type:Boolean,default:!0},onResize:Function}),setup(e){let{mergedClsPrefixRef:t,mergedRtlRef:n}=f(e),r=_(`InternalSelection`,n,t),i=z(null),a=z(null),o=z(null),s=z(null),c=z(null),l=z(null),u=z(null),d=z(null),p=z(null),m=z(null),h=z(!1),g=z(!1),v=z(!1),y=b(`InternalSelection`,`-internal-selection`,De,he,e,M(e,`clsPrefix`)),S=L(()=>e.clearable&&!e.disabled&&(v.value||e.active)),w=L(()=>e.selectedOption?e.renderTag?e.renderTag({option:e.selectedOption,handleClose:()=>{}}):e.renderLabel?e.renderLabel(e.selectedOption,!0):T(e.selectedOption[e.labelField],e.selectedOption,!0):e.placeholder),ee=L(()=>{let t=e.selectedOption;if(t)return t[e.labelField]}),te=L(()=>e.multiple?!!(Array.isArray(e.selectedOptions)&&e.selectedOptions.length):e.selectedOption!==null);function E(){var t;let{value:n}=i;if(n){let{value:r}=a;r&&(r.style.width=`${n.offsetWidth}px`,e.maxTagCount!==`responsive`&&((t=p.value)==null||t.sync({showAllItemsBeforeCalculate:!1})))}}function D(){let{value:e}=m;e&&(e.style.display=`none`)}function k(){let{value:e}=m;e&&(e.style.display=`inline-block`)}Y(M(e,`active`),e=>{e||D()}),Y(M(e,`pattern`),()=>{e.multiple&&x(E)});function A(t){let{onFocus:n}=e;n&&n(t)}function j(t){let{onBlur:n}=e;n&&n(t)}function ne(t){let{onDeleteOption:n}=e;n&&n(t)}function re(t){let{onClear:n}=e;n&&n(t)}function N(t){let{onPatternInput:n}=e;n&&n(t)}function ie(e){(!e.relatedTarget||!o.value?.contains(e.relatedTarget))&&A(e)}function P(e){o.value?.contains(e.relatedTarget)||j(e)}function F(e){re(e)}function I(){v.value=!0}function ae(){v.value=!1}function oe(t){!e.active||!e.filterable||t.target!==a.value&&t.preventDefault()}function R(e){ne(e)}let B=z(!1);function H(t){if(t.key===`Backspace`&&!B.value&&!e.pattern.length){let{selectedOptions:t}=e;t?.length&&R(t[t.length-1])}}let U=null;function W(t){let{value:n}=i;n&&(n.textContent=t.target.value,E()),e.ignoreComposition&&B.value?U=t:N(t)}function G(){B.value=!0}function K(){B.value=!1,e.ignoreComposition&&N(U),U=null}function ce(t){var n;g.value=!0,(n=e.onPatternFocus)==null||n.call(e,t)}function q(t){var n;g.value=!1,(n=e.onPatternBlur)==null||n.call(e,t)}function le(){var t,n;if(e.filterable)g.value=!1,(t=l.value)==null||t.blur(),(n=a.value)==null||n.blur();else if(e.multiple){let{value:e}=s;e?.blur()}else{let{value:e}=c;e?.blur()}}function ue(){var t,n,r;e.filterable?(g.value=!1,(t=l.value)==null||t.focus()):e.multiple?(n=s.value)==null||n.focus():(r=c.value)==null||r.focus()}function de(){let{value:e}=a;e&&(k(),e.focus())}function X(){let{value:e}=a;e&&e.blur()}function Z(e){let{value:t}=u;t&&t.setTextContent(`+${e}`)}function fe(){let{value:e}=d;return e}function pe(){return a.value}let me=null;function ge(){me!==null&&window.clearTimeout(me)}function _e(){e.active||(ge(),me=window.setTimeout(()=>{te.value&&(h.value=!0)},100))}function Q(){ge()}function ve(e){e||(ge(),h.value=!1)}Y(te,e=>{e||(h.value=!1)}),se(()=>{V(()=>{let t=l.value;t&&(e.disabled?t.removeAttribute(`tabindex`):t.tabIndex=g.value?-1:0)})}),be(o,e.onResize);let{inlineThemeDisabled:ye}=e,$=L(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:{fontWeight:r,borderRadius:i,color:a,placeholderColor:o,textColor:s,paddingSingle:c,paddingMultiple:l,caretColor:u,colorDisabled:d,textColorDisabled:f,placeholderColorDisabled:p,colorActive:m,boxShadowFocus:h,boxShadowActive:g,boxShadowHover:_,border:v,borderFocus:b,borderHover:x,borderActive:S,arrowColor:C,arrowColorDisabled:w,loadingColor:ee,colorActiveWarning:T,boxShadowFocusWarning:te,boxShadowActiveWarning:E,boxShadowHoverWarning:D,borderWarning:k,borderFocusWarning:A,borderHoverWarning:j,borderActiveWarning:ne,colorActiveError:re,boxShadowFocusError:M,boxShadowActiveError:N,boxShadowHoverError:ie,borderError:P,borderFocusError:F,borderHoverError:I,borderActiveError:ae,clearColor:L,clearColorHover:oe,clearColorPressed:R,clearSize:z,arrowSize:B,[J(`height`,t)]:V,[J(`fontSize`,t)]:H}}=y.value,U=O(c),W=O(l);return{"--n-bezier":n,"--n-border":v,"--n-border-active":S,"--n-border-focus":b,"--n-border-hover":x,"--n-border-radius":i,"--n-box-shadow-active":g,"--n-box-shadow-focus":h,"--n-box-shadow-hover":_,"--n-caret-color":u,"--n-color":a,"--n-color-active":m,"--n-color-disabled":d,"--n-font-size":H,"--n-height":V,"--n-padding-single-top":U.top,"--n-padding-multiple-top":W.top,"--n-padding-single-right":U.right,"--n-padding-multiple-right":W.right,"--n-padding-single-left":U.left,"--n-padding-multiple-left":W.left,"--n-padding-single-bottom":U.bottom,"--n-padding-multiple-bottom":W.bottom,"--n-placeholder-color":o,"--n-placeholder-color-disabled":p,"--n-text-color":s,"--n-text-color-disabled":f,"--n-arrow-color":C,"--n-arrow-color-disabled":w,"--n-loading-color":ee,"--n-color-active-warning":T,"--n-box-shadow-focus-warning":te,"--n-box-shadow-active-warning":E,"--n-box-shadow-hover-warning":D,"--n-border-warning":k,"--n-border-focus-warning":A,"--n-border-hover-warning":j,"--n-border-active-warning":ne,"--n-color-active-error":re,"--n-box-shadow-focus-error":M,"--n-box-shadow-active-error":N,"--n-box-shadow-hover-error":ie,"--n-border-error":P,"--n-border-focus-error":F,"--n-border-hover-error":I,"--n-border-active-error":ae,"--n-clear-size":z,"--n-clear-color":L,"--n-clear-color-hover":oe,"--n-clear-color-pressed":R,"--n-arrow-size":B,"--n-font-weight":r}}),xe=ye?C(`internal-selection`,L(()=>e.size[0]),$,e):void 0;return{mergedTheme:y,mergedClearable:S,mergedClsPrefix:t,rtlEnabled:r,patternInputFocused:g,filterablePlaceholder:w,label:ee,selected:te,showTagsPanel:h,isComposing:B,counterRef:u,counterWrapperRef:d,patternInputMirrorRef:i,patternInputRef:a,selfRef:o,multipleElRef:s,singleElRef:c,patternInputWrapperRef:l,overflowRef:p,inputTagElRef:m,handleMouseDown:oe,handleFocusin:ie,handleClear:F,handleMouseEnter:I,handleMouseLeave:ae,handleDeleteOption:R,handlePatternKeyDown:H,handlePatternInputInput:W,handlePatternInputBlur:q,handlePatternInputFocus:ce,handleMouseEnterCounter:_e,handleMouseLeaveCounter:Q,handleFocusout:P,handleCompositionEnd:K,handleCompositionStart:G,onPopoverUpdateShow:ve,focus:ue,focusInput:de,blur:le,blurInput:X,updateCounter:Z,getCounter:fe,getTail:pe,renderLabel:e.renderLabel,cssVars:ye?void 0:$,themeClass:xe?.themeClass,onRender:xe?.onRender}},render(){let{status:e,multiple:t,size:n,disabled:i,filterable:a,maxTagCount:o,bordered:s,clsPrefix:c,ellipsisTagPopoverProps:l,onRender:u,renderTag:f,renderLabel:p}=this;u?.();let h=o===`responsive`,g=typeof o==`number`,_=h||g,v=r(d,null,{default:()=>r(m,{clsPrefix:c,loading:this.loading,showArrow:this.showArrow,showClear:this.mergedClearable&&this.selected,onClear:this.handleClear},{default:()=>{var e;return(e=this.$slots).arrow?.call(e)}})}),y;if(t){let{labelField:e}=this,t=t=>r(`div`,{class:`${c}-base-selection-tag-wrapper`,key:t.value},f?f({option:t,handleClose:()=>{this.handleDeleteOption(t)}}):r(_e,{size:n,closable:!t.disabled,disabled:i,onClose:()=>{this.handleDeleteOption(t)},internalCloseIsButtonTag:!1,internalCloseFocusable:!1},{default:()=>p?p(t,!0):T(t[e],t,!0)})),s=()=>(g?this.selectedOptions.slice(0,o):this.selectedOptions).map(t),u=a?r(`div`,{class:`${c}-base-selection-input-tag`,ref:`inputTagElRef`,key:`__input-tag__`},r(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,tabindex:-1,disabled:i,value:this.pattern,autofocus:this.autofocus,class:`${c}-base-selection-input-tag__input`,onBlur:this.handlePatternInputBlur,onFocus:this.handlePatternInputFocus,onKeydown:this.handlePatternKeyDown,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),r(`span`,{ref:`patternInputMirrorRef`,class:`${c}-base-selection-input-tag__mirror`},this.pattern)):null,d=h?()=>r(`div`,{class:`${c}-base-selection-tag-wrapper`,ref:`counterWrapperRef`},r(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,onMouseleave:this.handleMouseLeaveCounter,disabled:i})):void 0,m;if(g){let e=this.selectedOptions.length-o;e>0&&(m=r(`div`,{class:`${c}-base-selection-tag-wrapper`,key:`__counter__`},r(_e,{size:n,ref:`counterRef`,onMouseenter:this.handleMouseEnterCounter,disabled:i},{default:()=>`+${e}`})))}let b=h?a?r(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,getTail:this.getTail,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:d,tail:()=>u}):r(ye,{ref:`overflowRef`,updateCounter:this.updateCounter,getCounter:this.getCounter,style:{width:`100%`,display:`flex`,overflow:`hidden`}},{default:s,counter:d}):g&&m?s().concat(m):s(),x=_?()=>r(`div`,{class:`${c}-base-selection-popover`},h?s():this.selectedOptions.map(t)):void 0,S=_?Object.assign({show:this.showTagsPanel,trigger:`hover`,overlap:!0,placement:`top`,width:`trigger`,onUpdateShow:this.onPopoverUpdateShow,theme:this.mergedTheme.peers.Popover,themeOverrides:this.mergedTheme.peerOverrides.Popover},l):null,C=!this.selected&&(!this.active||!this.pattern&&!this.isComposing)?r(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`},r(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):null,w=a?r(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-tags`},b,h?null:u,v):r(`div`,{ref:`multipleElRef`,class:`${c}-base-selection-tags`,tabindex:i?void 0:0},b,v);y=r(R,null,_?r(D,Object.assign({},S,{scrollable:!0,style:`max-height: calc(var(--v-target-height) * 6.6);`}),{trigger:()=>w,default:x}):w,C)}else if(a){let e=this.pattern||this.isComposing,t=this.active?!e:!this.selected,n=this.active?!1:this.selected;y=r(`div`,{ref:`patternInputWrapperRef`,class:`${c}-base-selection-label`,title:this.patternInputFocused?void 0:G(this.label)},r(`input`,Object.assign({},this.inputProps,{ref:`patternInputRef`,class:`${c}-base-selection-input`,value:this.active?this.pattern:``,placeholder:``,readonly:i,disabled:i,tabindex:-1,autofocus:this.autofocus,onFocus:this.handlePatternInputFocus,onBlur:this.handlePatternInputBlur,onInput:this.handlePatternInputInput,onCompositionstart:this.handleCompositionStart,onCompositionend:this.handleCompositionEnd})),n?r(`div`,{class:`${c}-base-selection-label__render-label ${c}-base-selection-overlay`,key:`input`},r(`div`,{class:`${c}-base-selection-overlay__wrapper`},f?f({option:this.selectedOption,handleClose:()=>{}}):p?p(this.selectedOption,!0):T(this.label,this.selectedOption,!0))):null,t?r(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},r(`div`,{class:`${c}-base-selection-overlay__wrapper`},this.filterablePlaceholder)):null,v)}else y=r(`div`,{ref:`singleElRef`,class:`${c}-base-selection-label`,tabindex:this.disabled?void 0:0},this.label===void 0?r(`div`,{class:`${c}-base-selection-placeholder ${c}-base-selection-overlay`,key:`placeholder`},r(`div`,{class:`${c}-base-selection-placeholder__inner`},this.placeholder)):r(`div`,{class:`${c}-base-selection-input`,title:G(this.label),key:`input`},r(`div`,{class:`${c}-base-selection-input__content`},f?f({option:this.selectedOption,handleClose:()=>{}}):p?p(this.selectedOption,!0):T(this.label,this.selectedOption,!0))),v);return r(`div`,{ref:`selfRef`,class:[`${c}-base-selection`,this.rtlEnabled&&`${c}-base-selection--rtl`,this.themeClass,e&&`${c}-base-selection--${e}-status`,{[`${c}-base-selection--active`]:this.active,[`${c}-base-selection--selected`]:this.selected||this.active&&this.pattern,[`${c}-base-selection--disabled`]:this.disabled,[`${c}-base-selection--multiple`]:this.multiple,[`${c}-base-selection--focus`]:this.focused}],style:this.cssVars,onClick:this.onClick,onMouseenter:this.handleMouseEnter,onMouseleave:this.handleMouseLeave,onKeydown:this.onKeydown,onFocusin:this.handleFocusin,onFocusout:this.handleFocusout,onMousedown:this.handleMouseDown},y,s?r(`div`,{class:`${c}-base-selection__border`}):null,s?r(`div`,{class:`${c}-base-selection__state-border`}):null)}});function ke(e){return e.type===`group`}function Ae(e){return e.type===`ignored`}function je(e,t){try{return!!(1+t.toString().toLowerCase().indexOf(e.trim().toLowerCase()))}catch{return!1}}function Me(e,t){return{getIsGroup:ke,getIgnored:Ae,getKey(t){return ke(t)?t.name||t.key||`key-required`:t[e]},getChildren(e){return e[t]}}}function Ne(e,t,n,r){if(!t)return e;function i(e){if(!Array.isArray(e))return[];let a=[];for(let o of e)if(ke(o)){let e=i(o[r]);e.length&&a.push(Object.assign({},o,{[r]:e}))}else if(Ae(o))continue;else t(n,o)&&a.push(o);return a}return i(e)}function Pe(e,t,n){let r=new Map;return e.forEach(e=>{ke(e)?e[n].forEach(e=>{r.set(e[t],e)}):r.set(e[t],e)}),r}var Fe=H([q(`select`,`
 z-index: auto;
 outline: none;
 width: 100%;
 position: relative;
 font-weight: var(--n-font-weight);
 `),q(`select-menu`,`
 margin: 4px 0;
 box-shadow: var(--n-menu-box-shadow);
 `,[t({originalTransition:`background-color .3s var(--n-bezier), box-shadow .3s var(--n-bezier)`})])]),Ie=Z({name:`Select`,props:Object.assign(Object.assign({},b.props),{to:c.propTo,bordered:{type:Boolean,default:void 0},clearable:Boolean,clearCreatedOptionsOnClear:{type:Boolean,default:!0},clearFilterAfterSelect:{type:Boolean,default:!0},options:{type:Array,default:()=>[]},defaultValue:{type:[String,Number,Array],default:null},keyboard:{type:Boolean,default:!0},value:[String,Number,Array],placeholder:String,menuProps:Object,multiple:Boolean,size:String,menuSize:{type:String},filterable:Boolean,disabled:{type:Boolean,default:void 0},remote:Boolean,loading:Boolean,filter:Function,placement:{type:String,default:`bottom-start`},widthMode:{type:String,default:`trigger`},tag:Boolean,onCreate:Function,fallbackOption:{type:[Function,Boolean],default:void 0},show:{type:Boolean,default:void 0},showArrow:{type:Boolean,default:!0},maxTagCount:[Number,String],ellipsisTagPopoverProps:Object,consistentMenuWidth:{type:Boolean,default:!0},virtualScroll:{type:Boolean,default:!0},labelField:{type:String,default:`label`},valueField:{type:String,default:`value`},childrenField:{type:String,default:`children`},renderLabel:Function,renderOption:Function,renderTag:Function,"onUpdate:value":[Function,Array],inputProps:Object,nodeProps:Function,ignoreComposition:{type:Boolean,default:!0},showOnFocus:Boolean,onUpdateValue:[Function,Array],onBlur:[Function,Array],onClear:[Function,Array],onFocus:[Function,Array],onScroll:[Function,Array],onSearch:[Function,Array],onUpdateShow:[Function,Array],"onUpdate:show":[Function,Array],displayDirective:{type:String,default:`show`},resetMenuOnOptionsChange:{type:Boolean,default:!0},status:String,showCheckmark:{type:Boolean,default:!0},scrollbarProps:Object,onChange:[Function,Array],items:Array}),slots:Object,setup(t){let{mergedClsPrefixRef:r,mergedBorderedRef:a,namespaceRef:o,inlineThemeDisabled:s,mergedComponentPropsRef:u}=f(t),d=b(`Select`,`-select`,Fe,ge,t,r),m=z(t.defaultValue),g=h(M(t,`value`),m),_=z(!1),v=z(``),y=l(t,[`items`,`options`]),x=z([]),S=z([]),w=L(()=>S.value.concat(x.value).concat(y.value)),ee=L(()=>{let{filter:e}=t;if(e)return e;let{labelField:n,valueField:r}=t;return(e,t)=>{if(!t)return!1;let i=t[n];if(typeof i==`string`)return je(e,i);let a=t[r];return typeof a==`string`?je(e,a):typeof a==`number`?je(e,String(a)):!1}}),T=L(()=>{if(t.remote)return y.value;{let{value:e}=w,{value:n}=v;return!n.length||!t.filterable?e:Ne(e,ee.value,n,t.childrenField)}}),te=L(()=>{let{valueField:e,childrenField:n}=t,r=Me(e,n);return le(T.value,r)}),E=L(()=>Pe(w.value,t.valueField,t.childrenField)),D=z(!1),O=h(M(t,`show`),D),k=z(null),A=z(null),j=z(null),{localeRef:ne}=n(`Select`),ie=L(()=>t.placeholder??ne.value.placeholder),P=[],F=z(new Map),I=L(()=>{let{fallbackOption:e}=t;if(e===void 0){let{labelField:e,valueField:n}=t;return t=>({[e]:String(t),[n]:t})}return e===!1?!1:t=>Object.assign(e(t),{value:t})});function ae(e){let n=t.remote,{value:r}=F,{value:i}=E,{value:a}=I,o=[];return e.forEach(e=>{if(i.has(e))o.push(i.get(e));else if(n&&r.has(e))o.push(r.get(e));else if(a){let t=a(e);t&&o.push(t)}}),o}let oe=L(()=>{if(t.multiple){let{value:e}=g;return Array.isArray(e)?ae(e):[]}return null}),R=L(()=>{let{value:e}=g;return!t.multiple&&!Array.isArray(e)?e===null?null:ae([e])[0]||null:null}),B=i(t,{mergedSize:e=>{let{size:n}=t;if(n)return n;let{mergedSize:r}=e||{};return r?.value?r.value:u?.value?.Select?.size||`medium`}}),{mergedSizeRef:V,mergedDisabledRef:H,mergedStatusRef:W}=B;function G(e,n){let{onChange:r,"onUpdate:value":i,onUpdateValue:a}=t,{nTriggerFormChange:o,nTriggerFormInput:s}=B;r&&U(r,e,n),a&&U(a,e,n),i&&U(i,e,n),m.value=e,o(),s()}function K(e){let{onBlur:n}=t,{nTriggerFormBlur:r}=B;n&&U(n,e),r()}function se(){let{onClear:e}=t;e&&U(e)}function ce(e){let{onFocus:n,showOnFocus:r}=t,{nTriggerFormFocus:i}=B;n&&U(n,e),i(),r&&X()}function q(e){let{onSearch:n}=t;n&&U(n,e)}function J(e){let{onScroll:n}=t;n&&U(n,e)}function ue(){var e;let{remote:n,multiple:r}=t;if(n){let{value:n}=F;if(r){let{valueField:r}=t;(e=oe.value)==null||e.forEach(e=>{n.set(e[r],e)})}else{let e=R.value;e&&n.set(e[t.valueField],e)}}}function de(e){let{onUpdateShow:n,"onUpdate:show":r}=t;n&&U(n,e),r&&U(r,e),D.value=e}function X(){H.value||(de(!0),D.value=!0,t.filterable&&Le())}function Z(){de(!1)}function fe(){v.value=``,S.value=P}let pe=z(!1);function me(){t.filterable&&(pe.value=!0)}function he(){t.filterable&&(pe.value=!1,O.value||fe())}function _e(){H.value||(O.value?t.filterable?Le():Z():X())}function Q(e){(j.value?.selfRef)?.contains(e.relatedTarget)||(_.value=!1,K(e),Z())}function ve(e){ce(e),_.value=!0}function ye(){_.value=!0}function be(e){k.value?.$el.contains(e.relatedTarget)||(_.value=!1,K(e),Z())}function $(){var e;(e=k.value)==null||e.focus(),Z()}function xe(t){O.value&&(k.value?.$el.contains(e(t))||Z())}function Se(e){if(!Array.isArray(e))return[];if(I.value)return Array.from(e);{let{remote:n}=t,{value:r}=E;if(n){let{value:t}=F;return e.filter(e=>r.has(e)||t.has(e))}else return e.filter(e=>r.has(e))}}function Ce(e){we(e.rawNode)}function we(e){if(H.value)return;let{tag:n,remote:r,clearFilterAfterSelect:i,valueField:a}=t;if(n&&!r){let{value:e}=S,t=e[0]||null;if(t){let e=x.value;e.length?e.push(t):x.value=[t],S.value=P}}if(r&&F.value.set(e[a],e),t.multiple){let t=Se(g.value),o=t.findIndex(t=>t===e[a]);if(~o){if(t.splice(o,1),n&&!r){let t=Te(e[a]);~t&&(x.value.splice(t,1),i&&(v.value=``))}}else t.push(e[a]),i&&(v.value=``);G(t,ae(t))}else{if(n&&!r){let t=Te(e[a]);~t?x.value=[x.value[t]]:x.value=P}Ie(),Z(),G(e[a],e)}}function Te(e){return x.value.findIndex(n=>n[t.valueField]===e)}function Ee(e){O.value||X();let{value:n}=e.target;v.value=n;let{tag:r,remote:i}=t;if(q(n),r&&!i){if(!n){S.value=P;return}let{onCreate:e}=t,r=e?e(n):{[t.labelField]:n,[t.valueField]:n},{valueField:i,labelField:a}=t;y.value.some(e=>e[i]===r[i]||e[a]===r[a])||x.value.some(e=>e[i]===r[i]||e[a]===r[a])?S.value=P:S.value=[r]}}function De(e){e.stopPropagation();let{multiple:n,tag:r,remote:i,clearCreatedOptionsOnClear:a}=t;!n&&t.filterable&&Z(),r&&!i&&a&&(x.value=P),se(),n?G([],[]):G(null,null)}function Oe(e){!N(e,`action`)&&!N(e,`empty`)&&!N(e,`header`)&&e.preventDefault()}function ke(e){J(e)}function Ae(e){var n,r,i;if(!t.keyboard){e.preventDefault();return}switch(e.key){case` `:if(t.filterable)break;e.preventDefault();case`Enter`:if(!k.value?.isComposing){if(O.value){let e=j.value?.getPendingTmNode();e?Ce(e):t.filterable||(Z(),Ie())}else if(X(),t.tag&&pe.value){let e=S.value[0];if(e){let n=e[t.valueField],{value:r}=g;t.multiple&&Array.isArray(r)&&r.includes(n)||we(e)}}}e.preventDefault();break;case`ArrowUp`:if(e.preventDefault(),t.loading)return;O.value&&((n=j.value)==null||n.prev());break;case`ArrowDown`:if(e.preventDefault(),t.loading)return;O.value?(r=j.value)==null||r.next():X();break;case`Escape`:O.value&&(re(e),Z()),(i=k.value)==null||i.focus();break}}function Ie(){var e;(e=k.value)==null||e.focus()}function Le(){var e;(e=k.value)==null||e.focusInput()}function Re(){var e;O.value&&((e=A.value)==null||e.syncPosition())}ue(),Y(M(t,`options`),ue);let ze={focus:()=>{var e;(e=k.value)==null||e.focus()},focusInput:()=>{var e;(e=k.value)==null||e.focusInput()},blur:()=>{var e;(e=k.value)==null||e.blur()},blurInput:()=>{var e;(e=k.value)==null||e.blurInput()}},Be=L(()=>{let{self:{menuBoxShadow:e}}=d.value;return{"--n-menu-box-shadow":e}}),Ve=s?C(`select`,void 0,Be,t):void 0;return Object.assign(Object.assign({},ze),{mergedStatus:W,mergedClsPrefix:r,mergedBordered:a,namespace:o,treeMate:te,isMounted:p(),triggerRef:k,menuRef:j,pattern:v,uncontrolledShow:D,mergedShow:O,adjustedTo:c(t),uncontrolledValue:m,mergedValue:g,followerRef:A,localizedPlaceholder:ie,selectedOption:R,selectedOptions:oe,mergedSize:V,mergedDisabled:H,focused:_,activeWithoutMenuOpen:pe,inlineThemeDisabled:s,onTriggerInputFocus:me,onTriggerInputBlur:he,handleTriggerOrMenuResize:Re,handleMenuFocus:ye,handleMenuBlur:be,handleMenuTabOut:$,handleTriggerClick:_e,handleToggle:Ce,handleDeleteOption:we,handlePatternInput:Ee,handleClear:De,handleTriggerBlur:Q,handleTriggerFocus:ve,handleKeydown:Ae,handleMenuAfterLeave:fe,handleMenuClickOutside:xe,handleMenuScroll:ke,handleMenuKeydown:Ae,handleMenuMousedown:Oe,mergedTheme:d,cssVars:s?void 0:Be,themeClass:Ve?.themeClass,onRender:Ve?.onRender})},render(){return r(`div`,{class:`${this.mergedClsPrefix}-select`},r(y,null,{default:()=>[r(X,null,{default:()=>r(Oe,{ref:`triggerRef`,inlineThemeDisabled:this.inlineThemeDisabled,status:this.mergedStatus,inputProps:this.inputProps,clsPrefix:this.mergedClsPrefix,showArrow:this.showArrow,maxTagCount:this.maxTagCount,ellipsisTagPopoverProps:this.ellipsisTagPopoverProps,bordered:this.mergedBordered,active:this.activeWithoutMenuOpen||this.mergedShow,pattern:this.pattern,placeholder:this.localizedPlaceholder,selectedOption:this.selectedOption,selectedOptions:this.selectedOptions,multiple:this.multiple,renderTag:this.renderTag,renderLabel:this.renderLabel,filterable:this.filterable,clearable:this.clearable,disabled:this.mergedDisabled,size:this.mergedSize,theme:this.mergedTheme.peers.InternalSelection,labelField:this.labelField,valueField:this.valueField,themeOverrides:this.mergedTheme.peerOverrides.InternalSelection,loading:this.loading,focused:this.focused,onClick:this.handleTriggerClick,onDeleteOption:this.handleDeleteOption,onPatternInput:this.handlePatternInput,onClear:this.handleClear,onBlur:this.handleTriggerBlur,onFocus:this.handleTriggerFocus,onKeydown:this.handleKeydown,onPatternBlur:this.onTriggerInputBlur,onPatternFocus:this.onTriggerInputFocus,onResize:this.handleTriggerOrMenuResize,ignoreComposition:this.ignoreComposition},{arrow:()=>{var e;return[(e=this.$slots).arrow?.call(e)]}})}),r(B,{ref:`followerRef`,show:this.mergedShow,to:this.adjustedTo,teleportDisabled:this.adjustedTo===c.tdkey,containerClass:this.namespace,width:this.consistentMenuWidth?`target`:void 0,minWidth:`target`,placement:this.placement},{default:()=>r(a,{name:`fade-in-scale-up-transition`,appear:this.isMounted,onAfterLeave:this.handleMenuAfterLeave},{default:()=>{var e;return this.mergedShow||this.displayDirective===`show`?((e=this.onRender)==null||e.call(this),F(r(Ee,Object.assign({},this.menuProps,{ref:`menuRef`,onResize:this.handleTriggerOrMenuResize,inlineThemeDisabled:this.inlineThemeDisabled,virtualScroll:this.consistentMenuWidth&&this.virtualScroll,class:[`${this.mergedClsPrefix}-select-menu`,this.themeClass,this.menuProps?.class],clsPrefix:this.mergedClsPrefix,focusable:!0,labelField:this.labelField,valueField:this.valueField,autoPending:!0,nodeProps:this.nodeProps,theme:this.mergedTheme.peers.InternalSelectMenu,themeOverrides:this.mergedTheme.peerOverrides.InternalSelectMenu,treeMate:this.treeMate,multiple:this.multiple,size:this.menuSize,renderOption:this.renderOption,renderLabel:this.renderLabel,value:this.mergedValue,style:[this.menuProps?.style,this.cssVars],onToggle:this.handleToggle,onScroll:this.handleMenuScroll,onFocus:this.handleMenuFocus,onBlur:this.handleMenuBlur,onKeydown:this.handleMenuKeydown,onTabOut:this.handleMenuTabOut,onMousedown:this.handleMenuMousedown,show:this.mergedShow,showCheckmark:this.showCheckmark,resetMenuOnOptionsChange:this.resetMenuOnOptionsChange,scrollbarProps:this.scrollbarProps}),{empty:()=>{var e;return[(e=this.$slots).empty?.call(e)]},header:()=>{var e;return[(e=this.$slots).header?.call(e)]},action:()=>{var e;return[(e=this.$slots).action?.call(e)]}}),this.displayDirective===`show`?[[s,this.mergedShow],[ne,this.handleMenuClickOutside,void 0,{capture:!0}]]:[[ne,this.handleMenuClickOutside,void 0,{capture:!0}]])):null}})})]}))}});export{$ as i,Me as n,Ee as r,Ie as t};